package main

import (
	"pentestplatform/attack"
	"pentestplatform/web"
)

func main() {
	attack.LoadPlugin()
	web.Run()
}
