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
	attacker, err := attack.GetAttacker("tp5_0_23RCE")
	if err != nil{
		logger.Red.Printf("%s run failed", "sql")
	}
	attacker.SetOptions("{\"target\":\"http://127.0.0.1:8080/\",\"cmd\":\"\"}")
	_, isVuln := attacker.Exploit()
	fmt.Println(isVuln)
}
