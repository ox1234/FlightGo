package main

import (
	"fmt"
	"pentestplatform/gather"
)

func main(){
	vtScanner := gather.NewVtScanner()
	vtScanner.Set("xidian.edu.cn")
	vtScanner.DoGather()
	jsondata, _ := vtScanner.Report()
	fmt.Println(jsondata)
}
