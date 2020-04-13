package main

import "pentestplatform/gather"

func main(){
	dirScanner := gather.NewDirScanner("php72.f1ight.top", "php")
	dirScanner.DoGather()
}