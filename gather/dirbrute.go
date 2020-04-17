package gather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"pentestplatform/util"
)

type dirScanner struct {
	Target      string
	wordlist    []string
	Dirs        []dirDesc
	concurrency int
}

type dirDesc struct {
	Url string
	StatusCode int
	ContentLength int
}

func NewDirScanner() *dirScanner {
	return &dirScanner{
		concurrency: 100,
	}
}

func (d *dirScanner) Set(v ...interface{}){
	target := v[0].(string)
	dirType := v[1].(string)
	dictName := fmt.Sprintf("dict/dir-%s.txt", dirType)
	d.wordlist = util.ReadFile(dictName)
	d.Target = target
}

func (d *dirScanner) Report() (string, error){
	jsondata, err := json.Marshal(d)
	if err != nil{
		logger.Red.Println(err)
		return "", err
	}
	return string(jsondata), nil

}

func (d *dirScanner) DoGather(){
	tracker := make(chan bool)
	dirnames := make(chan string)
	for i:=0; i<d.concurrency; i++{
		go d.worker(tracker, dirnames)
	}

	for _, dirname := range d.wordlist{
		dirnames <- dirname
	}
	close(dirnames)

	for i:=0; i<d.concurrency; i++{
		<- tracker
	}
	logger.Green.Println("dirbrute complete")
}

func (d *dirScanner) worker(tracker chan bool, dirnames chan string){
	for dirname := range dirnames{
		d.fetch(dirname)
	}

	var empty bool
	tracker <- empty
}

func (d *dirScanner) fetch(dirname string){
	if d.checkAlive(){
		url := fmt.Sprintf("http://%s%s", d.Target, dirname)
		resp, err := http.Get(url)
		if err != nil{
			logger.Red.Fatal(err)
			return
		}
		statusCode := resp.StatusCode
		if statusCode != 404{
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			contentLength := len(body)
			d.Dirs = append(d.Dirs, dirDesc{
				Url:           url,
				StatusCode:    statusCode,
				ContentLength: contentLength,
			})
		}
	}
}

func (d *dirScanner) checkAlive() bool{
	_, err := http.Get(fmt.Sprintf("http://%s/", d.Target))
	if err != nil{
		logger.Red.Fatalf("%s is not alive ", d.Target)
		return false
	}
	return true
}


