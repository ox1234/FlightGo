package gather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ullaakut/nmap"
	"github.com/dean2021/go-masscan"
	"github.com/gosuri/uilive"
	"pentestplatform/logger"
	"strconv"
	"strings"
	"time"
)

type masScanner struct {
	concurrency int
	ip2scan []string
	IpMap map[string][]Port
}


type Port struct {
	Port string
	Service string
	State string
}

func NewPortScanner() *masScanner{
	return &masScanner{
		concurrency: 2,
		IpMap: make(map[string][]Port),
	}
}

func (m *masScanner) Set(v ...interface{}){
	m.ip2scan = append(m.ip2scan, v[0].(string))
}

func (m *masScanner) DoGather(){
	writer := uilive.New()
	writer.Start()
	tracker := make(chan bool)
	ips := make(chan string)
	for i:=0; i<m.concurrency; i++{
		go m.worker(tracker, ips)
	}

	for i, ip := range m.ip2scan{
		fmt.Fprintf(writer, "端口扫描进度：%d/%d\n", i, len(m.ip2scan))
		time.Sleep(time.Millisecond * 5)
		ips <- ip
	}

	close(ips)
	writer.Stop()
	for i:=0; i<m.concurrency; i++{
		<- tracker
	}
	logger.Green.Println("port scan complete")
}

func (m *masScanner) Report()(string, error){
	jsondata, err := json.Marshal(m.IpMap)
	if err != nil{
		return "", err
	}
	return string(jsondata), nil
}

func (m *masScanner) worker(tracker chan bool, ips chan string){
	for ip := range ips{
		m.doScan(ip)
	}
	var empty bool
	tracker <- empty
}

func (m *masScanner) doScan(ip string){
	openPorts := m.doMasScan(ip)
	m.doNmapScan(ip, openPorts)
}

func (m *masScanner) doNmapScan(ip string, openPorts []string){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var nmapScanner *nmap.Scanner
	if len(openPorts) > 0{
		logger.Blue.Println("do nmap and masscan")
		scanner, err := nmap.NewScanner(
			nmap.WithTargets(ip),
			nmap.WithPorts(strings.Join(openPorts,",")),
			nmap.WithContext(ctx),
			nmap.WithSkipHostDiscovery(),
			nmap.WithConnectScan(),
			nmap.WithServiceInfo(),
		)
		if err != nil{
			logger.Red.Println("fail to create nmap scanner ", err)
			return
		}
		nmapScanner = scanner
	}else{
		logger.Blue.Println("do normal nmap scan")
		scanner, err := nmap.NewScanner(
			nmap.WithTargets(ip),
			nmap.WithContext(ctx),
			nmap.WithSkipHostDiscovery(),
			nmap.WithConnectScan(),
			nmap.WithServiceInfo(),
		)
		if err != nil{
			logger.Red.Println("fail to create nmap scanner ", err)
			return
		}
		nmapScanner = scanner
	}
	result, warnings, err := nmapScanner.Run()
	if err != nil {
		logger.Red.Println("unable to run nmap scan: %v", err)
		return
	}

	if warnings != nil {
		logger.Blue.Printf("Warnings: \n %v\n", warnings)
	}

	for _, host := range result.Hosts{
		if len(host.Ports) == 0 || len(host.Addresses) == 0{
			fmt.Println(ip, " port is not open")
			continue
		}
		for _, port := range host.Ports{
			m.IpMap[ip] = append(m.IpMap[ip], Port{
				Port:    strconv.Itoa(int(port.ID)),
				Service: port.Service.Name,
				State:   port.State.State,
			})
		}
	}
}

func (m *masScanner) doMasScan(ip string)[]string{
	mass := masscan.New()
	mass.SetSystemPath("/usr/local/bin/masscan")
	mass.SetPorts("0-65535")
	mass.SetRanges(ip)
	mass.SetRate("1000")
	err := mass.Run()
	if err != nil{
		logger.Red.Println("scan failed ", err)
		return []string{}
	}
	var openPorts []string
	results, err := mass.Parse()
	if err != nil{
		logger.Red.Println("parse failed ", err)
		return []string{}
	}
	for _, host := range results{
		if len(host.Ports) > 0{
			for _, port := range host.Ports{
				openPorts = append(openPorts, port.Portid)
			}
		}
	}
	return openPorts
}
