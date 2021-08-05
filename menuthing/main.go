package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"gopkg.in/yaml.v3"
)

func main() {
	systray.Run(onReady, onExit)
}

type Config struct {
	Menu MenuConfig
}

type MenuConfig struct {
	Icon  string
	Items []MenuConfigItem
}

type MenuConfigItem struct {
	Label string
	Type  string
	Open  string
	Exec  string
	TZ    string
}

func (mci *MenuConfigItem) GetLabel() string {
	if mci.Label != "" {
		return mci.Label
	}
	if mci.Open != "" {
		return mci.Open
	}
	if mci.Exec != "" {
		return mci.Exec
	}
	if mci.Type != "" {
		return mci.Type
	}
	return ""
}

func (mci *MenuConfigItem) GetTimeFormat() string {
	if mci.Label != "" {
		return mci.Label
	}
	return time.RFC3339
}

func onReady() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Open(u.HomeDir + "/.menuthing.yaml")
	if err != nil {
		log.Fatal(err)
	}
	d := yaml.NewDecoder(fp)

	var config Config
	d.Decode(&config)
	for _, item := range config.Menu.Items {
		switch item.Type {
		case "separator", "--":
			systray.AddSeparator()
		case "clock":
			addClock(item)
		default:
			addMenuAction(item.GetLabel(), "", action(item))
		}
	}

	// fmt.Printf("CONFIG>> %#v\n", config)

	setIcon(config.Menu.Icon)

	systray.AddSeparator()
	addMenuAction("Exit", "", func(_ *systray.MenuItem) { onExit() })
}

func onExit() {
	fmt.Println("cheers")
	os.Exit(0)
}

func addClock(item MenuConfigItem) {
	m := systray.AddMenuItem("", "")
	location := time.Local
	if item.TZ != "" {
		var err error
		location, err = time.LoadLocation(item.TZ)
		if err != nil {
			m.SetTitle(err.Error())
			return
		}
	}
	go func() {
		for {
			m.SetTitle(time.Now().In(location).Format(item.GetTimeFormat()))
			time.Sleep(5 * time.Second)
		}
	}()
}

func addMenuAction(label, tooltip string, f func(*systray.MenuItem)) *systray.MenuItem {
	m := systray.AddMenuItem(label, tooltip)
	go func() {
		for {
			<-m.ClickedCh
			go f(m)
		}
	}()
	return m
}

func action(item MenuConfigItem) func(*systray.MenuItem) {
	return func(m *systray.MenuItem) {
		fmt.Printf(">> %#v\n", item)
		if item.Open != "" {
			open.Run(item.Open)
		}
		if item.Exec != "" {
			ff := strings.Fields(item.Exec)
			c := exec.Command(ff[0])
			if len(ff) > 1 {
				c.Args = ff
			}
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			fmt.Printf("running %v\n", c)

			if err := c.Run(); err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("cool")
			}
		}
	}
}

func setIcon(name string) {
	buf, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		// convert to ico if it's not already
		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".png" {
			buf = toIco(png.Decode, buf)
		} else if ext == ".jpg" || ext == ".jpeg" {
			buf = toIco(jpeg.Decode, buf)
		}
	}

	systray.SetIcon(buf)
}

func toIco(f func(r io.Reader) (image.Image, error), inBits []byte) []byte {
	img, err := f(bytes.NewReader(inBits))
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := ico.Encode(buf, img); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
