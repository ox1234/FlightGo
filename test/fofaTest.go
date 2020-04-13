package main

import (
	"pentestplatform/gather"
)

func main(){
	fofaScanner := gather.NewFofaScanner()
	fofaScanner.Set("xidian.edu.cn")
	fofaScanner.DoGather()
}
