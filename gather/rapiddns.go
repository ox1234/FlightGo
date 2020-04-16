package gather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"regexp"
)

type rapidDnsScanner struct {
	Domain string
	RapidDomainSet []subDomain
}

const searchURL string = "http://rapiddns.io/subdomain/%s?full=1&down=1"

func NewRapidDnsScanner() *rapidDnsScanner{
	return &rapidDnsScanner{}
}


func (r *rapidDnsScanner) Set(v ...interface{}){
	r.Domain = v[0].(string)
}

func (r *rapidDnsScanner) DoGather(){
	finalUrl := fmt.Sprintf(searchURL, r.Domain)
	resp, err := http.Get(finalUrl)
	if err != nil{
		logger.Red.Println(err)
		return
	}
	bodyReader := resp.Body
	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil{
		logger.Red.Println(err)
		return
	}
	bodyString := string(bodyBytes)
	r.extractDomain(bodyString)
	logger.Green.Println("rapiddns scan over")
}

func (r *rapidDnsScanner) Report()(string, error){
	jsondata, err := json.Marshal(r)
	if err != nil{
		logger.Red.Println(err)
		return "", err
	}
	return string(jsondata), nil
}

func (r *rapidDnsScanner) extractDomain(rdhtml string){
	regexpPattern, err := regexp.Compile(`<th scope="row ">\d*</th>
<td><a href=".*" target="_blank">(.*)</a></td>
<td><a href=".*" target="_blank" title=".*">(.*)</a>
</td>`)
	if err != nil{
		logger.Red.Println(err)
		return
	}
	allMatch := regexpPattern.FindAllStringSubmatch(rdhtml, -1)
	for _, perMatch := range allMatch{
		ip := perMatch[2]
		domain := perMatch[1]
		r.RapidDomainSet = append(r.RapidDomainSet, subDomain{
			IPAddress: ip,
			HostName:  domain,
		})
	}
}
