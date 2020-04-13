package gather

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/miekg/dns"
	"math/rand"
	"pentestplatform/logger"
	"pentestplatform/util"
	"time"
)

type SubDomainScanner struct {
	rootDomain string
	SubDomains []subDomain
	concurrency int
	wordlist []string
	dnsServer []string
}

type subDomain struct {
	IPAddress string
	HostName string
}

func NewSubDomainScanner() *SubDomainScanner{
	return &SubDomainScanner{
		wordlist: util.ReadFile("dict/subtest.txt"),
		concurrency: 100,
		dnsServer: []string{
			"114.114.114.114:53",
			"114.114.115.115:53",
			"223.5.5.5:53",
			"223.6.6.6:53",
			"180.76.76.76:53",
			"119.29.29.29:53",
			"182.254.116.116:53",
			"1.2.4.8:53",
			"8.8.8.8:53",
			"8.8.4.4:53",
			"1.1.1.1:53",
			"1.0.0.1:53",
		},
	}
}

func (s *SubDomainScanner) Set(v ...interface{}){
	s.rootDomain = v[0].(string)
}

func(s *SubDomainScanner) Report() (string, error){
	subDomainResult := struct {
		RootDomain string
		SubDomains []subDomain
	}{
		RootDomain: s.rootDomain,
		SubDomains: s.SubDomains,
	}
	jsondata, err := json.Marshal(subDomainResult)
	if err != nil{
		logger.Red.Fatal(err)
		return "", err
	}
	return string(jsondata), nil
}

func (s *SubDomainScanner) DoGather(){
	writer := uilive.New()
	writer.Start()
	fqdns := make(chan string, s.concurrency)
	tracker := make(chan bool)
	for i:=0; i < s.concurrency; i++{
		go s.worker(tracker, fqdns, s.dnsServer[rand.Intn(len(s.dnsServer))])
	}

	for index, sub := range s.wordlist{
		fmt.Fprintf(writer, "子域名扫描进度：%d/%d\n", index, len(s.wordlist))
		time.Sleep(time.Millisecond * 5)
		fqdns <- fmt.Sprintf("%s.%s", sub, s.rootDomain)
	}

	writer.Stop()
	close(fqdns)
	for i:=0; i<s.concurrency; i++{
		<- tracker
	}

	logger.Green.Println("enum complete")
}

func (s *SubDomainScanner) Lookup(fqdn, serverAddr string){
	var cfqdns = fqdn
	for{
		cname, err := lookupCNAME(fqdn, serverAddr)
		if err == nil && len(cname) > 0{
			cfqdns = cname[0]
			continue
		}
		ips, err := lookupA(cfqdns, serverAddr)
		if err != nil{
			break
		}
		for _, ip := range ips{
			//logger.Blue.Println("hostname: " + fqdn)
			//logger.Blue.Println("ip: " + ip)
			if !s.hasScanned(fqdn, ip){
				s.SubDomains = append(s.SubDomains, subDomain{
					IPAddress: ip,
					HostName:  fqdn,
				})
			}
		}
		break
	}
}

func (s *SubDomainScanner) hasScanned(domain, ip string) bool{
	for _, sub := range s.SubDomains{
		if sub.HostName == domain && sub.IPAddress == ip{
			return true
		}
	}
	return false
}

func (s *SubDomainScanner) worker(tracker chan bool, fqdns chan string, serverAddr string){
	for fqdn := range fqdns{
		s.Lookup(fqdn, serverAddr)
	}

	var empty bool
	tracker <- empty
}

func lookupA(fqdn, serverAddr string)([]string, error){
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil{
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips, errors.New("no answer")
	}
	for _, answer := range in.Answer{
		if a, ok := answer.(*dns.A); ok{
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func lookupCNAME(fqdn, serverAddr string)([]string, error){
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil{
		return fqdns, err
	}
	if len(in.Answer) < 1{
		return fqdns, errors.New("no answer")
	}
	for _, answer := range in.Answer{
		if c, ok := answer.(*dns.CNAME); ok{
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}



