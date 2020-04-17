package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"pentestplatform/attack"
	"pentestplatform/logger"
)

var attackExploiter = attack.NewAttackExploiter()


func ShowPayload(context *gin.Context){
	pluginJson, err := json.Marshal(attackExploiter)
	if err != nil{
		logger.Red.Println(err)
	}

	context.String(200, string(pluginJson))
}

func DoExploit(context *gin.Context){
	attackerName := context.PostForm("attName")
	attackerOptions := context.PostForm("options")
	attacker, ok := attackExploiter.PluginMap[attackerName]
	if !ok{
		context.String(404, "no such attacker")
		return
	}
	attacker.Attacker.SetOptions(attackerOptions)
	body, isVuln := attacker.Attacker.Exploit()
	if isVuln{
		context.String(200, "this site is vuln")
	}else{
		context.String(200, fmt.Sprintf("please check %s", body))
	}
}