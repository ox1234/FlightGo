package gather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pentestplatform/logger"
	"pentestplatform/util"
)

type dirTarget struct {
	target string
	wordlist []string
	dirs []dirDesc
	concurrency int
}

type dirDesc struct {
	Url string
	StatusCode int
	ContentLength int
}

func NewDirScanner() *dirTarget{
	return &dirTarget{
		concurrency: 100,
	}
}

func (d *dirTarget) Set(v ...interface{}){
	target := v[0].(string)
	dirType := v[1].(string)
	dictName := fmt.Sprintf("dict/dir-%s.txt", dirType)
	d.wordlist = util.ReadFile(dictName)
	d.target = target
}

func (d *dirTarget) Report() (string, error){
	dirTarget := struct {
		Target string
		Dirs   []dirDesc
	}{
		Target: d.target,
		Dirs: d.dirs,
	}
	jsondata, err := json.Marshal(dirTarget)
	if err != nil{
		logger.Red.Fatal(err)
		return "", err
	}
	return string(jsondata), nil

}

func (d *dirTarget) DoGather(){
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

func (d *dirTarget) worker(tracker chan bool, dirnames chan string){
	for dirname := range dirnames{
		d.fetch(dirname)
	}

	var empty bool
	tracker <- empty
}

func (d *dirTarget) fetch(dirname string){
	if d.checkAlive(){
		url := fmt.Sprintf("http://%s%s", d.target, dirname)
		fmt.Println(url)
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
			d.dirs = append(d.dirs, dirDesc{
				Url:           url,
				StatusCode:    statusCode,
				ContentLength: contentLength,
			})
		}
	}
}

func (d *dirTarget) checkAlive() bool{
	_, err := http.Get(fmt.Sprintf("http://%s/", d.target))
	if err != nil{
		logger.Red.Fatalf("%s is not alive ", d.target)
		return false
	}
	return true
}


