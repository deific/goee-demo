package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//将request转发给 http://127.0.0.1:2003
func helloHandler(w http.ResponseWriter, r *http.Request) {

	//trueServer := "http://59.255.86.249:8081/"
	//trueServer := "https://app.yideb.com/"
	trueServer := "http://www.baidu.com/"
	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("proxy to:", url)
	proxy := httputil.NewSingleHostReverseProxy(url)
	// Update the headers to allow for SSL redirection
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}

func main() {
	fmt.Println("------HttpProxy---------")
	listenPort := ":8081"
	log.Println("Start listen port:", listenPort)
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
