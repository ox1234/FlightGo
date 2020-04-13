package gather

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uilive"
	"net"
	"pentestplatform/logger"
	"pentestplatform/util"
	"strconv"
	"strings"
	"time"
)

type portScanner struct {
	ipDescs []ipDesc
	ipConcurrency int
	portList []string
	portConcurrency int
}

type ipDesc struct {
	ip string
	ports []portDesc
}

type portDesc struct {
	port int
	service string
}

func NewPortScanner() *portScanner{
	return &portScanner{
		ipConcurrency:   10,
		portList:        util.ReadFile("dict/Top100ports.txt"),
		portConcurrency: 10,
	}
}

func (p *portScanner) Set(v ...interface{}){
	ips := v[0].([]string)
	for _, ip := range ips{
		if !p.hasScanned(ip){
			ipDesc := ipDesc{
				ip:    ip,
			}
			p.ipDescs = append(p.ipDescs, ipDesc)
		}
	}
}

func (p *portScanner) hasScanned(ip string) bool{
	for _, ipdesc := range p.ipDescs{
		if ipdesc.ip == ip{
			return true
		}
	}
	return false
}

func (p *portScanner) DoGather(){
	writer := uilive.New()
	writer.Start()
	ipTracker := make(chan bool)
	ips := make(chan *ipDesc)

	for i:=0; i<p.ipConcurrency; i++{
		go p.ipWorker(ipTracker, ips)
	}
	for i:=0; i<len(p.ipDescs); i++{
		fmt.Fprintf(writer, "端口扫描进度：%d/%d\n", i, len(p.ipDescs))
		time.Sleep(time.Millisecond * 5)
		ips <- &p.ipDescs[i]
	}

	close(ips)
	writer.Stop()
	for i:=0; i<p.ipConcurrency; i++{
		<- ipTracker
	}
	logger.Green.Println("port scan complete")
}

func (p *portScanner) Report() (string, error){
	type portDesc struct {
		Port int
		Service string
	}
	type ipDesc struct {
		Ip string
		Ports []portDesc
	}
	var result []ipDesc
	for _, ipdesc := range p.ipDescs{
		ip := ipdesc.ip
		var portarr []portDesc
		for _, portdesc := range ipdesc.ports{
			service := portdesc.service
			port := portdesc.port
			portarr = append(portarr, portDesc{
				Port:    port,
				Service: service,
			})
		}
		result = append(result, ipDesc{
			Ip:    ip,
			Ports: portarr,
		})
	}
	jsondata, err := json.Marshal(result)
	if err != nil{
		logger.Red.Fatal(err)
		return "", err
	}
	return string(jsondata), nil
}

func (p *portScanner) ipWorker(ipTrack chan bool, ips chan *ipDesc){
	for ip := range ips{
		p.scanPort(ip)
	}

	var empty bool
	ipTrack <- empty
}

func (p *portScanner) scanPort(ipdesc *ipDesc){
	ports := make(chan string, p.portConcurrency)
	portTracker := make(chan bool)
	for i:=0; i<p.portConcurrency; i++{
		go ipdesc.portWorker(portTracker, ports)
	}

	for _, port := range p.portList{
		ports <- port
	}

	close(ports)
	for i:=0; i<p.portConcurrency; i++{
		<- portTracker
	}
}

func (i *ipDesc) portWorker(portTracker chan bool, ports chan string){
	for port := range ports{
		portarr := strings.Split(port, " ")
		intPort, _ := strconv.Atoi(portarr[0])
		i.checkAlive(intPort, portarr[1])
	}

	var empty bool
	portTracker <- empty
}


func (i *ipDesc) checkAlive(port int, desc string){
	addr := fmt.Sprintf("%s:%d", i.ip, port)
	_, err := net.DialTimeout("tcp", addr, 3 * time.Second)
	if err == nil{
		i.ports = append(i.ports, portDesc{
			port:    port,
			service: desc,
		})
	}
}



