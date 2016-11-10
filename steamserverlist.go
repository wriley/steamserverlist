package main

import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
	"bytes"
	//"unsafe"
	"net/url"
	"flag"
	"strings"
)

type SteamServerList struct {
	Response struct {
		Servers []struct {
			Addr string `json:"addr"`
			Gameport int `json:"gameport"`
			Steamid string `json:"steamid"`
			Name string `json:"name"`
			Appid int `json:"appid"`
			Gamedir string `json:"gamedir"`
			Version string `json:"version"`
			Product string `json:"product"`
			Region int `json:"region"`
			Players int `json:"players"`
			MaxPlayers int `json:"max_players"`
			Bots int `json:"bots"`
			Map string `json:"map"`
			Secure bool `json:"secure"`
			Dedicated bool `json:"dedicated"`
			Os string `json:"os"`
			Gametype string `json:"gametype"`
		} `json:"servers"`
	} `json:"response"`
}

func main() {
	KeyPtr := flag.String("key", "", "Steam API key")
	LimitPtr := flag.Int("limit", 5000, "Limit search results")
	FilterPtr := flag.String("filter", "", "filter string")
	
	flag.Parse()
	
	if *KeyPtr == "" {
		fmt.Printf("-key is required\n")
		return
	}
	
	SafeKey := url.QueryEscape(*KeyPtr)
	SafeFilter := url.QueryEscape(*FilterPtr)
	
    url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?limit=%d&key=%s", *LimitPtr, SafeKey)
	
	if len(SafeFilter) > 0 {
		url = fmt.Sprintf("%s&filter=%s", url, SafeFilter)
	}
	
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var serverlist SteamServerList
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	b := buf.Bytes()
	//s := *(*string)(unsafe.Pointer(&b))
	
	//fmt.Println(s)
	
	if err := json.Unmarshal(b, &serverlist); err != nil {
		log.Fatal(err)
	}
	
	for _, server := range serverlist.Response.Servers {
		tokens := strings.Split(server.Addr, ":")
		fmt.Printf("%s %s\n", tokens[0], tokens[1])
	}
}
