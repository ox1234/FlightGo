package deleted

// 没钱，暂时不使用fofa

import (
	"fmt"
	"github.com/fofapro/fofa-go/fofa"
	"pentestplatform/logger"
)

const FOFA_EMAIL string = "923435274@qq.com"
const FOFA_KEY string = "4e8f541604217ec0ccb808ff6134ef7f"

type fofaScanner struct {
	Domain string

}

func NewFofaScanner() *fofaScanner{
	return &fofaScanner{}
}

func (f *fofaScanner) Set(v ...interface{}){
	f.Domain = v[0].(string)
}

func (f *fofaScanner) DoGather(){
	clt := fofa.NewFofaClient([]byte(FOFA_EMAIL), []byte(FOFA_KEY))
	if clt == nil{
		logger.Red.Println("fofa client construct error!")
	}
	fofaQuery := fmt.Sprintf("domain=\"%s\"", f.Domain)
	result, err := clt.QueryAsJSON(1, []byte(fofaQuery))
	if err != nil{
		logger.Red.Println(err)
	}
	fmt.Printf("%s\n", result)
}

