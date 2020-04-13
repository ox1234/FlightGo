package main

import (
	"encoding/json"
	"log"
)

type ipDesc struct {
	Ip    string
	Ports []portDesc
}

type portDesc struct {
	Port    int
	Service string
}

type test struct {
	name string `json:"jafksdjlfk"`
}

func main() {
	test2 := test{name:"aksdjflk"}
	jsondata, err := json.Marshal(test2)
	if err != nil{

	}
	log.Printf(string(jsondata))
}
