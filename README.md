# light-frame

> A digital picture frame project written in Go and hobbled together with a collection of cron jobs

## build for Raspberry Pi Zero W
- `env GOOS=linux GOARCH=arm GOARM=6 go build -o ../executable/light-frame-app-pi -C app -v`

If you have Go installed on your system, simply run `go run app/main.go` from the project root.
