package attack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"pentestplatform/logger"
	plugin2 "plugin"
	"strings"
)

type Attacker interface {
	/*
	arg1: 攻击的目标
	argn: 额外的参数，根据需要添加
	 */
	Exploit(v ...interface{})(string, error)
	Describe()string
}

const pluginDir string = "plugin/"

var pluginMap map[string]Attacker

func LoadPlugin(){
	pluginMap = make(map[string]Attacker)
	files, _ := ioutil.ReadDir(pluginDir)
	for _, file := range files{
		if strings.HasSuffix(file.Name(), ".so"){
			soFile := file.Name()
			pluginName := strings.Split(soFile, ".")[0]
			fattacker, err := findPlugin(soFile)
			if err == nil{
				pluginMap[pluginName] = fattacker
			}
		}
	}
	for k, _ := range pluginMap{
		logger.Blue.Printf("%s is loaded", k)
	}
}

func GetAttacker(module string)(Attacker, error){
	attacker, ok := pluginMap[module]
	if !ok{
		logger.Red.Printf("%s is not loaded or has error", module)
		return nil, errors.New("module is not loaded or has error")
	}
	return attacker, nil
}

func ShowPlugin()string{
	var pluginNameArr []string
	for k, _ := range pluginMap{
		pluginNameArr = append(pluginNameArr, k)
	}
	pluginJson, err := json.Marshal(pluginNameArr)
	if err != nil{
		logger.Red.Println(err)
	}
	return string(pluginJson)
}

func findPlugin(soFile string)(Attacker, error){
	soPath := pluginDir + soFile
	plugin, err := plugin2.Open(soPath)
	if err != nil{
		logger.Red.Printf("plugin %s load failed|%s", soFile, err.Error())
		return nil, err
	}
	symAttacker, err := plugin.Lookup("FAttacker")
	if err != nil{
		logger.Red.Printf("plugin %s load failed|%s", soFile, err.Error())
		return nil, err
	}
	var attacker Attacker
	attacker, ok := symAttacker.(Attacker)
	if !ok{
		logger.Red.Printf("plugin %s load failed|uexpected type from module")
		return nil, errors.New("uexpected type from module")
	}
	return attacker, nil
}