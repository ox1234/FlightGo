package subdomain

import (
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

type Record struct {
	Domain string
	Type   string
	Target string
	IP     []string
}

var (
	Queries    = make(chan string)
	Records    = make(chan *Record)
	serverAddr string
	rateLimit  int
	client     *dns.Conn
	noQueries  = make(chan bool)
)

func Configure(dnsServer string, rate int) {
	serverAddr = dnsServer
	rateLimit = rate
	client, _ = dns.DialTimeout("udp", serverAddr, 3*time.Second)
	go send()
	go receive()
}

func send() {
	delay := time.Second / time.Nanosecond / time.Duration(rateLimit)
	for {
		select {
		case domain := <-Queries:
			time.Sleep(delay)
			m := dns.Msg{}
			m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
			if err := client.WriteMsg(&m); err != nil {
				log.Println(err)
				continue
			}
		case <-time.After(2 * time.Second):
			close(noQueries)
			return
		}
	}
}

func receive() {
	for {
		select {
		case <-noQueries:
			close(Records)
			return
		default:
		}
		var msg *dns.Msg
		var err error
		if err = client.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			continue
		}
		if msg, err = client.ReadMsg(); err != nil || len(msg.Answer) == 0 {
			continue
		}
		domain := strings.TrimSuffix(msg.Question[0].Name, ".")
		if record := newRecord(domain, msg.Answer); record != nil {
			Records <- record
		}
	}
}

func newRecord(domain string, response []dns.RR) *Record {
	record := Record{Domain: domain}
	switch firstAnswer := response[0].(type) {
	case *dns.CNAME:
		record.Type = "CNAME"
		record.Target = strings.TrimSuffix(firstAnswer.Target, ".")
		response = response[1:]
	case *dns.A:
		record.Type = "A"
	default:
		return nil
	}

	for _, ans := range response {
		if a, ok := ans.(*dns.A); ok {
			record.IP = append(record.IP, a.A.String())
		}
	}

	return &record
}
