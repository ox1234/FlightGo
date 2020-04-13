package main

import (
	"pentestplatform/gather"
	"pentestplatform/logger"
	"sync"
)

func main(){
	subDomainScanner := gather.NewSubDomainScanner("xidian.edu.cn")
	subDomainScanner.DoGather()
	subjson, err := subDomainScanner.Report()
	if err == nil{
		logger.Blue.Printf("%s", subjson)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		var iplist []string
		for _, subdomain := range subDomainScanner.SubDomains{
			iplist = append(iplist, subdomain.IPAddress)
		}
		portScanner := gather.NewPortScanner(iplist)
		portScanner.DoGather()
		json, err := portScanner.Report()
		if err == nil{
			logger.Blue.Printf("%s", json)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		for _, subDomain := range subDomainScanner.SubDomains{
			basicScanner := gather.NewBasicScanner(subDomain.HostName, subDomain.IPAddress)
			basicScanner.DoGather()
			jsondata, err := basicScanner.Report()
			if err != nil{
				logger.Red.Fatal(err)
			}
			logger.Blue.Printf("%s", jsondata)
		}
		wg.Done()
	}()
	wg.Wait()
}