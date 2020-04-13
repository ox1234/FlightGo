package gather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"strings"
)

type basicInfo struct {
	HostName string
	Ip string
	WebServer string
	ClickJackingProtection bool
	ContentSecurityPolicy bool
	XContentTypeOptions bool
	StrictTransportSecurity bool
}

func NewBasicScanner() *basicInfo{
	return &basicInfo{}
}

func (b *basicInfo) Set(v ...interface{}){
	hostname := v[0].(string)
	ip := v[1].(string)
	b.HostName = hostname
	b.Ip = ip
}

func (b *basicInfo) DoGather(){
	b.doGet()
	logger.Green.Println("basic info complete")
}

func (b *basicInfo) Report() (string, error){
	jsondata, err := json.Marshal(b)
	if err != nil{
		logger.Red.Fatal(err)
		return "", err
	}
	return string(jsondata), err

}

func (b *basicInfo) doGet(){
	url := fmt.Sprintf("http://%s/", b.HostName)
	resp, err := http.Get(url)
	if err != nil{
		logger.Red.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if len(resp.Header["X-Frame-Options"]) > 0 {
		b.ClickJackingProtection = true
	}
	if len(resp.Header["Content-Security-Policy"]) > 0 ||
		strings.Contains(string(body), `http-equiv="Content-Security-Policy"`) {
		b.ContentSecurityPolicy = true
	}
	if len(resp.Header["X-Content-Type-Options"]) > 0 {
		b.XContentTypeOptions = true
	}
	if len(resp.Header["Strict-Transport-Secruity"]) > 0 {
		b.StrictTransportSecurity = true
	}
	if len(resp.Header["Server"]) > 0{
		b.WebServer = resp.Header["Server"][0]
	}
}
