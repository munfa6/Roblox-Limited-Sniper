package data

import (
	"bufio"
	"crypto/tls"
	"math/rand"
	"os"
	"time"
)

type Proxys struct {
	Proxys []string
	Used   map[string]bool
	Conn   *tls.Conn
}

func (Proxy *Proxys) GetProxys(uselist bool, list []string) {
	Proxy.Proxys = []string{}
	if uselist {
		Proxy.Proxys = append(Proxy.Proxys, list...)
	} else {
		file, err := os.Open("proxys.txt")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				Proxy.Proxys = append(Proxy.Proxys, scanner.Text())
			}
		}
	}
}

func (Proxy *Proxys) CompRand() string {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	return Proxy.Proxys[rand.Intn(len(Proxy.Proxys))]
}

func (Proxy *Proxys) Setup() {
	Proxy.Used = make(map[string]bool)
	for _, proxy := range Proxy.Proxys {
		Proxy.Used[proxy] = false
	}
}

func (Proxy *Proxys) RandProxy() string {
	for _, proxy := range Proxy.Proxys {
		if !Proxy.Used[proxy] {
			Proxy.Used[proxy] = true
			return proxy
		}
	}

	Proxy.Setup()

	return ""
}
