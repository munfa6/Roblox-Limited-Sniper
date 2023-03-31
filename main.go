package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	D     []Items
	Robux int64
)

func init() {
	data.Con.LoadState()
	data.Proxy.GetProxys(false, nil)
	data.Proxy.Setup()
	Clear()

	go refresh()

	if file_name := "proxys.txt"; data.CheckForValidFile(file_name) {
		os.Create(file_name)
	}

	GetUpdatedItemList()
	go func() {
		for {
			time.Sleep(time.Second)
			GetUpdatedItemList()
		}
	}()
	Robux = data.GetRobux()
	go func() {
		for {
			time.Sleep(time.Minute)
			Robux = data.GetRobux()
		}
	}()

	fmt.Print(`
 ██▓    
▓██▒    
▒██░    
▒██░    
░██████▒
░ ▒░▓  ░
░ ░ ▒  ░
  ░ ░   
    ░  ░
`)

}

func Client(PS []string) (Client *http.Client) {
	if len(PS) > 2 {
		Client = (&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: PS[0] + ":" + PS[1], User: url.UserPassword(PS[2], PS[3])})}})
	} else {
		Client = (&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: PS[0] + ":" + PS[1]})}})
	}
	return
}

func main() {

	type Item struct {
		Success    bool            `json:"success"`
		Activities [][]interface{} `json:"activities"`
	}

	for {
		go func() {
			req, _ := http.NewRequest("GET", "https://www.rolimons.com/api/activity2", nil)
			PS := strings.Split(data.Proxy.CompRand(), ":")
			if resp, err := Client(PS).Do(req); err == nil && resp.StatusCode == 200 {
				resp, _ := io.ReadAll(resp.Body)
				var F Item
				json.Unmarshal(resp, &F)
				var Cached map[string]Items = make(map[string]Items)
				for _, f := range F.Activities {
					if _, ok := Cached[f[2].(string)]; !ok {
						for _, k := range D {
							if k.ID == f[2].(string) {
								Cached[f[2].(string)] = k
								break
							}
						}
					}
				}
			Exit:
				for _, ok := range Cached {
					if Percent(ok.AvgPrice, ok.BestPrice) >= data.Con.Percentage && Robux > int64(ok.BestPrice) {
						for _, d := range data.Con.PreferredIds {
							if ok.ID == fmt.Sprintf("%v", d) {
								buyitem(ok.ID, int(ok.BestPrice))
								break Exit
							}
						}
						buyitem(ok.ID, int(ok.BestPrice))
					}
				}
			}
		}()
		time.Sleep(100 * time.Millisecond)
	}
}

func Percent(avg, best float64) float64 {
	return math.RoundToEven(((avg - best) / avg) * 100)
}

func buyitem(id string, price int) {
	r1, _ := http.NewRequest("GET", fmt.Sprintf("https://economy.roblox.com/v1/assets/%v/resellers?cursor=&limit=100", id), nil)
	r1.AddCookie(&http.Cookie{
		Name:  ".ROBLOSECURITY",
		Value: data.Con.ROBLOSECURITY,
	})
	resp, _ := http.DefaultClient.Do(r1)
	if resp.StatusCode == 200 {
		f, _ := io.ReadAll(resp.Body)
		var Info data.SellerInfo
		json.Unmarshal(f, &Info)
		if len(Info.Data) > 0 && Percent((float64(price)-(0.30*float64(price))), float64(Info.Data[1].Price)) >= data.Con.BuyersProfit {
			req, _ := http.NewRequest("POST", "https://economy.roblox.com/v1/purchases/products/"+id, bytes.NewBuffer([]byte(fmt.Sprintf(`{"expectedCurrency":1,"expectedPrice":%v,"expectedSellerId":%v,"userAssetId":%v}`, price, Info.Data[0].Seller.ID, Info.Data[0].UserAssetID))))
			req.Header.Add("x-csrf-token", data.Csrf)
			req.Header.Add("cookie", ".ROBLOSECURITY="+data.Con.ROBLOSECURITY)
			req.Header.Add("Content-Type", "application/json")
			if resp, err := http.DefaultClient.Do(req); err == nil && resp.StatusCode == 200 {
				b, _ := io.ReadAll(resp.Body)
				var P data.Purchased
				json.Unmarshal(b, &P)
				if P.Purchased {
					fmt.Printf("Succesfully purchased %v.\n", P.ProductID)
				} else {
					fmt.Println(P.ErrorMsg)
				}
			} else {
				fmt.Println(resp.StatusCode)
				f, _ := io.ReadAll(resp.Body)
				fmt.Println(string(f))
			}
		}
		//
	} else {
		fmt.Println(resp.StatusCode)
		f, _ := io.ReadAll(resp.Body)
		fmt.Println(string(f))
	}
}

type Items struct {
	ID        string
	Name      string
	AvgPrice  float64
	BestPrice float64
}

func GetUpdatedItemList() {
	req, _ := http.NewRequest("GET", "https://www.rolimons.com/deals", nil)
	PS := strings.Split(data.Proxy.CompRand(), ":")
	if resp, err := Client(PS).Do(req); err == nil && resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		var M map[string]interface{}
		json.Unmarshal([]byte(strings.Split(strings.Split(string(body), `var item_details =`)[1], `;`)[0]), &M)
		var N []Items
		for key, data := range M {
			N = append(N, Items{
				ID:        key,
				Name:      data.([]interface{})[0].(string),
				BestPrice: data.([]interface{})[1].(float64),
				AvgPrice:  data.([]interface{})[2].(float64),
			})
		}
		D = N
	}
}

func refresh() {
	for {
		if resp, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", nil); err == nil {
			resp.Header.Add("cookie", ".ROBLOSECURITY="+data.Con.ROBLOSECURITY)
			r, _ := http.DefaultClient.Do(resp)
			data.Csrf = r.Header.Get("x-csrf-token")
		}
		time.Sleep(10 * time.Minute)
	}
}

func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
