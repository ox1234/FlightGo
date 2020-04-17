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
	Exploit()(string, bool)
	GetDesc()string
	GetTitle()string
	GetVulnType()string
	GetOptions()string
	SetOptions(optionJson string)
	IsVulnerable(v ...interface{})bool
}

type AttackPlugin struct {
	Title    string
	Attacker Attacker `json:"-"`
	Desc     string
	VulnType string
	Options  string
}

type AttackExploiter struct {
	PluginMap map[string]AttackPlugin
	PluginNum int
}

const pluginDir string = "plugin/"

var pluginMap map[string]AttackPlugin

func NewAttackExploiter() *AttackExploiter{
	return &AttackExploiter{
		PluginMap: pluginMap,
		PluginNum: len(pluginMap),
	}
}

func LoadPlugin(){
	pluginMap = make(map[string]AttackPlugin)
	files, _ := ioutil.ReadDir(pluginDir)
	for _, file := range files{
		if strings.HasSuffix(file.Name(), ".so"){
			soFile := file.Name()
			pluginName := strings.Split(soFile, ".")[0]
			fattacker, err := findPlugin(soFile)
			if err == nil{
				pluginMap[pluginName] = AttackPlugin{
					Title:    fattacker.GetTitle(),
					Attacker: fattacker,
					Desc:     fattacker.GetDesc(),
					VulnType: fattacker.GetVulnType(),
					Options:  fattacker.GetOptions(),
				}
			}
		}
	}
	for k, _ := range pluginMap{
		logger.Blue.Printf("%s is loaded", k)
	}
}

func GetAttacker(module string)(Attacker, error){
	attackerPlugin, ok := pluginMap[module]
	if !ok{
		logger.Red.Printf("%s is not loaded or has error", module)
		return nil, errors.New("module is not loaded or has error")
	}
	return attackerPlugin.Attacker, nil
}

func ShowPlugin()string{
	pluginJson, err := json.Marshal(pluginMap)
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