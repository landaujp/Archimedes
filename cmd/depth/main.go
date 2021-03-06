package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/landaujp/archimedes/depth"
)

//go:generate go-bindata config/config.toml

type Config struct {
	DB struct {
		Host     string
		Database string
		User     string
		Password string
		Port     int
	}
}
type Exchange interface {
	GetDepth() string
	GetBid() int
	GetAsk() int
	SetJson(*simplejson.Json)
}

func main() {

	dboption := redis.DialDatabase(0)
	con, err := redis.Dial("tcp", "127.0.0.1:6379", dboption)
	if err != nil {
		// handle error
	}
	defer con.Close()

	flag.Parse()
	argument := flag.Args()[0]

	var ex Exchange
	var url string
	switch argument {
	case "bitflyer":
		ex = &depth.Bitflyer{}
		url = "https://api.bitflyer.jp/v1/getboard"
	case "coincheck":
		ex = &depth.Coincheck{}
		url = "https://coincheck.com/api/order_books"
	case "zaif":
		ex = &depth.Zaif{}
		url = "https://api.zaif.jp/api/1/depth/btc_jpy"
	case "bitbank":
		ex = &depth.Bitbank{}
		url = "https://public.bitbank.cc/btc_jpy/depth"
	case "kraken":
		ex = &depth.Kraken{}
		url = "https://api.kraken.com/0/public/Depth?pair=XBTJPY"
	case "quoine":
		ex = &depth.Quoine{}
		url = "https://api.quoine.com/products/5/price_levels"
	case "btcbox":
		ex = &depth.Btcbox{}
		url = "https://www.btcbox.co.jp/api/v1/depth/"
	case "fisco":
		ex = &depth.Fisco{}
		url = "https://api.fcce.jp/api/1/depth/btc_jpy"
	default:
		fmt.Println("There is no exchanges...")
		return
	}

	var Etag string
	for {
		time.Sleep(5 * time.Second) // 2秒待つ

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("if-none-match", Etag)
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		if resp.StatusCode != 200 {
			continue
		}

		if val, ok := resp.Header["Etag"]; ok {
			Etag = val[0]
		}

		body, _ := ioutil.ReadAll(resp.Body)
		jsonObj, _ := simplejson.NewJson(body)
		ex.SetJson(jsonObj)
		jsonString := ex.GetDepth()

		con.Do("SET", argument+":depth", jsonString)
		con.Do("SET", argument+":bid", ex.GetBid())
		con.Do("SET", argument+":ask", ex.GetAsk())
		con.Do("EXPIRE", argument+":bid", 60)
		con.Do("EXPIRE", argument+":ask", 60)

		resp.Body.Close()
	}
}
