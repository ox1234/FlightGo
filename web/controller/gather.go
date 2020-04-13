package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"pentestplatform/gather"
	"strings"
	"sync"
)

var siteInfo []string

var subDomainScanner = gather.NewSubDomainScanner()
var portScanner = gather.NewPortScanner()
var dirScanner = gather.NewDirScanner()
var basicScanner = gather.NewBasicScanner()
var vtScanner = gather.NewVtScanner()

func SubDomain(context *gin.Context){
	jsondata, err := subDomainScanner.Report()
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, jsondata)
}

func VtDomain(context *gin.Context){
	jsondata, err := vtScanner.Report()
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, jsondata)
}

func PortScan(context *gin.Context){
	json, err := portScanner.Report()
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, json)
}

func DirScan(context *gin.Context){
	domain := strings.TrimSpace(context.PostForm("domain"))
	dirType := strings.TrimSpace(context.PostForm("type"))
	dirScanner.Set(domain, dirType)
	dirScanner.DoGather()
	json, err := dirScanner.Report()
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, json)
}

func BasicScan(context *gin.Context){
	resultjson, err := json.Marshal(siteInfo)
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, string(resultjson))
}

func Start(context *gin.Context){
	wg := sync.WaitGroup{}
	wg.Add(1)
	domain := strings.TrimSpace(context.PostForm("domain"))
	go func() {
		subDomainScanner.Set(domain)
		subDomainScanner.DoGather()
		wg.Done()
	}()
	wg.Wait()

	go func(){
		vtScanner.Set(domain)
		vtScanner.DoGather()
	}()

	go func() {
		var iplist []string
		for _, subdomain := range subDomainScanner.SubDomains{
			iplist = append(iplist, subdomain.IPAddress)
		}
		portScanner.Set(iplist)
		portScanner.DoGather()
	}()

	go func() {
		for _, subDomain := range subDomainScanner.SubDomains{
			hostname := subDomain.HostName
			ip := subDomain.IPAddress
			basicScanner.Set(hostname, ip)
			basicScanner.DoGather()
			json, err := basicScanner.Report()
			if err != nil{
				log.Fatal(err)
			}
			siteInfo = append(siteInfo, json)
		}
	}()
	context.String(200, "start scan, please wait")
}
