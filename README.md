# menuthing

```sh
go install ./menuthing
```

Example config goes in `~/.menuthing.yaml`

```yml
menu:
  icon: /home/reed/Pictures/mini-darling.png
  items:
    - label: kittens
      open: https://kittens.nz

    - type: ----

    - type: clock
      label: 3:04pm Brisbane
      tz: Australia/Brisbane

    - type: clock
      label: 3:04pm
      open: https://kittens.nz

    - type: clock
      label: 3:04pm Los Angeles
      tz: America/Los_Angeles

    - type: clock
      label: 3:04pm Edmonton
      tz: America/Edmonton

    - type: clock
      label: 3:04pm Indianapolis
      tz: America/Indianapolis

    - type: ----

    - open: https://pkg.go.dev/

    - exec: pi-down

    - exec: terminator -m -x 5mins
      label: 5 Mins

    - open: https://github.com/reedwade/menuthing
      label: github.com/reedwade/menuthing

```
