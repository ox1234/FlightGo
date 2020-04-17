package main

import (
	"pentestplatform/attack"
	"pentestplatform/web"
)

func main(){
	/*
	开启插件功能
	 */
	attack.LoadPlugin()
	web.Run()
}
