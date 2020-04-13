package main

import (
	"fmt"
	"pentestplatform/attack"
	"pentestplatform/logger"
)

func main(){
	attack.LoadPlugin()
	pluginJson := attack.ShowPlugin()
	fmt.Println(pluginJson)
	attacker, err := attack.GetAttacker("sql")
	if err != nil{
		logger.Red.Printf("%s run failed", "sql")
	}
	attacker.Exploit("fklasjdlkfjasdkf")
}
