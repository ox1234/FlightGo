package subdomain

import (
	"encoding/json"
	"fmt"
	"pentestplatform/logger"
	"pentestplatform/util"
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
		wordlist: util.ReadFile("dict/subnames.txt"),
		concurrency: 10,
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
	Configure("8.8.8.8:53", 10000)
	go func() {
		for _, sub := range s.wordlist{
			domain := fmt.Sprintf("%s.%s", sub, s.rootDomain)
			Queries <- domain
		}
	}()
	for record := range Records{
		s.SubDomains = append(s.SubDomains, subDomain{
			IPAddress: record.IP[0],
			HostName:  record.Domain,
		})
	}

	logger.Green.Println("enum complete")
}
