package util

import (
	"bufio"
	"os"
	"pentestplatform/logger"
	"strings"
)

func ReadFile(filename string)(data []string){
	fd, err := os.Open(filename)
	if err != nil{
		logger.Red.Fatal(err)
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan(){
		data = append(data, strings.TrimSpace(scanner.Text()))
	}
	return data
}
