package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"pentestplatform/gather"
	"pentestplatform/logger"
	"strings"
	"sync"
)

var siteInfo []string

var subDomainScanner = gather.NewSubDomainScanner()
var portScanner = gather.NewPortScanner()
var dirScanner = gather.NewDirScanner()
var basicScanner = gather.NewBasicScanner()
var vtScanner = gather.NewVtScanner()
var rapidDnsScanner = gather.NewRapidDnsScanner()


var allDomain = make(map[string]string)

func SubDomain(context *gin.Context){
	jsondata, err := subDomainScanner.Report()
	if err != nil{
		logger.Red.Println(err)
		context.String(500, "error")
	}
	context.String(200, jsondata)
}

func VtDomain(context *gin.Context){
	jsondata, err := vtScanner.Report()
	if err != nil{
		logger.Red.Println(err)
		context.String(500, "error")
	}
	context.String(200, jsondata)
}

func RapidDnsDomain(context *gin.Context){
	jsondata, err := rapidDnsScanner.Report()
	if err != nil{
		logger.Red.Println(err)
		context.String(500, "error")
	}
	context.String(200, jsondata)
}

func AllDomain(context *gin.Context){
	jsondata, err := json.Marshal(allDomain)
	if err != nil{
		logger.Red.Println(err)
		context.String(500, "error")
	}
	context.String(200, string(jsondata))
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
	jsondata, err := basicScanner.Report()
	if err != nil{
		log.Fatal(err)
	}
	context.String(200, jsondata)
}



func Start(context *gin.Context){
	wg := sync.WaitGroup{}
	domain := strings.TrimSpace(context.PostForm("domain"))
	go func() {
		/*
		进行子域名收集
		1. 子域名爆破
		2. vt上搜集子域名
		3. rapidDns上收集子域名
		 */
		wg.Add(1)
		go func() {
			subDomainScanner.Set(domain)
			subDomainScanner.DoGather()
			wg.Done()
		}()

		wg.Add(1)
		go func(){
			vtScanner.Set(domain)
			vtScanner.DoGather()
			wg.Done()
		}()
		wg.Add(1)
		go func(){
			rapidDnsScanner.Set(domain)
			rapidDnsScanner.DoGather()
			wg.Done()
		}()
		wg.Wait()
		collectDomain()
		/*
		因为网络带宽和性能平静瓶颈，暂时不开启端口扫描
		 */

		//go func() {
		//	iplist := make(map[string]bool)
		//	for _, ipaddress := range allDomain{
		//		if ipaddress != ""{
		//			iplist[ipaddress] = true
		//		}
		//	}
		//
		//	for ip := range iplist{
		//		portScanner.Set(ip)
		//	}
		//	portScanner.DoGather()
		//}()

		go func() {
			for subdomain, ipaddress := range allDomain{
				hostname := subdomain
				ip := ipaddress
				basicScanner.Set(hostname, ip)
			}
			basicScanner.DoGather()
		}()
	}()
	context.String(200, "start scan, please wait")
}

func collectDomain(){
	for _, sub := range rapidDnsScanner.RapidDomainSet{
		allDomain[sub.HostName] = sub.IPAddress
	}
	for _, sub := range vtScanner.VtDomainSet{
		allDomain[sub.HostName] = sub.IPAddress
	}
	for _, sub := range subDomainScanner.SubDomains{
		allDomain[sub.HostName] = sub.IPAddress
	}
}