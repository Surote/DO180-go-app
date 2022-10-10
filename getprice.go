package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var redis_pass = os.Getenv("REDIS_PASS")

var client = redis.NewClient(&redis.Options{
	Addr:     "surote-redis-headless.back-linkerd:6379",
	Password: redis_pass,
	DB:       0,
})


func main() {
	srv := &http.Server{
		Addr:         ":8020",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/hello-trigger", HelloServerTrigger)
	http.HandleFunc("/redis", GetRedis)
	http.HandleFunc("/", GetPrice)

	log.Fatal(srv.ListenAndServe())
}

func HelloServer(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Your path is , %s!", r.URL.Path[1:])
	fmt.Println(r.URL)
}

func HelloServerTrigger(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Your trigger path is , %s!", r.URL.Path[1:])
	fmt.Println(r.URL)
}

func GetRedis(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	val, err := client.Get("c").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	fmt.Fprintf(w, "%s", string(val))
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	var t map[string]interface{}
	resp, err := http.Get("https://data.messari.io/api/v1/assets/" + r.URL.Path[1:] + "/metrics")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &t)
	data, ok := t["data"].(map[string]interface{})
	if ok {
		fmt.Println("ok, 1 ")
	}
	market_price, ok := data["market_data"].(map[string]interface{})
	if ok {
		fmt.Println("ok, 2 ")
	}
	b, _ := json.Marshal(market_price)
	fmt.Println(resp.Status)
	fmt.Fprintf(w, "%s", string(b))
}
