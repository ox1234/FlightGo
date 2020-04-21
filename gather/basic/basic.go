package basic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"regexp"
	"strings"
)

type basicInfo struct {
	HostName string
	Ip string
	WebServer string
	Title string
	ClickJackingProtection bool
	ContentSecurityPolicy bool
	XContentTypeOptions bool
	StrictTransportSecurity bool
}

type basicScanner struct {
	ScanList []basicInfo
	concurrency int
}

func NewBasicScanner() *basicScanner{
	return &basicScanner{
		concurrency: 100,
	}
}

func (b *basicScanner) Set(v ...interface{}){
	hostname := v[0].(string)
	ip := v[1].(string)
	b.ScanList = append(b.ScanList, basicInfo{
		HostName: hostname,
		Ip:       ip,
	})
}

func (b *basicScanner) DoGather(){
	tracker := make(chan bool)
	sites := make(chan *basicInfo)

	for i:=0; i<b.concurrency; i++{
		go b.worker(tracker, sites)
	}

	for i:=0; i<len(b.ScanList); i++{
		sites <- &b.ScanList[i]
	}

	close(sites)
	for i:=0; i<b.concurrency; i++{
		<- tracker
	}
	logger.Green.Println("basic info complete")
}

func (b *basicScanner) Report() (string, error){
	jsondata, err := json.Marshal(b.ScanList)
	if err != nil{
		logger.Red.Fatal(err)
		return "", err
	}
	return string(jsondata), err
}

func (b *basicScanner) doGet(basicInfo *basicInfo){
	url := fmt.Sprintf("http://%s/", basicInfo.HostName)
	resp, err := http.Get(url)
	if err != nil{
		logger.Red.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile(`<title>(\s*.*\s*)</title>`)
	results := r.FindStringSubmatch(string(body))
	if len(results) == 2{
		basicInfo.Title = strings.TrimSpace(results[1])
	}else{
		basicInfo.Title = ""
	}

	defer resp.Body.Close()
	if len(resp.Header["X-Frame-Options"]) > 0 {
		basicInfo.ClickJackingProtection = true
	}else{
		basicInfo.ClickJackingProtection = false
	}
	if len(resp.Header["Content-Security-Policy"]) > 0 ||
		strings.Contains(string(body), `http-equiv="Content-Security-Policy"`) {
		basicInfo.ContentSecurityPolicy = true
	}else{
		basicInfo.ContentSecurityPolicy = false
	}
	if len(resp.Header["X-Content-Type-Options"]) > 0 {
		basicInfo.XContentTypeOptions = true
	}else{
		basicInfo.XContentTypeOptions = false
	}
	if len(resp.Header["Strict-Transport-Secruity"]) > 0 {
		basicInfo.StrictTransportSecurity = true
	}else{
		basicInfo.StrictTransportSecurity = false
	}
	if len(resp.Header["Server"]) > 0{
		basicInfo.WebServer = resp.Header["Server"][0]
	}else{
		basicInfo.WebServer = ""
	}
}

func (b basicScanner) worker(tracker chan bool, sites chan *basicInfo){
	for site := range sites{
		b.doGet(site)
	}

	var empty bool
	tracker <- empty
}
