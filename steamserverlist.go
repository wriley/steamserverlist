package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "bytes"
    "net/url"
    "flag"
    "strings"
)

// nested struct to hold json data
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
    // Get command line arguments
    KeyPtr := flag.String("key", "", "Steam API key")
    LimitPtr := flag.Int("limit", 5000, "Limit search results")
    FilterPtr := flag.String("filter", "", "filter string")
    flag.Parse()

    // Steam API key is required
    if *KeyPtr == "" {
        fmt.Printf("-key is required\n")
        return
    }

    // Escape user input and build API url
    SafeKey := url.QueryEscape(*KeyPtr)
    SafeFilter := url.QueryEscape(*FilterPtr)
    url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?limit=%d&key=%s", *LimitPtr, SafeKey)

    // If a filter was supplied on command line then add it to url
    if len(SafeFilter) > 0 {
        url = fmt.Sprintf("%s&filter=%s", url, SafeFilter)
    }

    // Setup http request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal("NewRequest: ", err)
        return
    }

    // Use http.Client to send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal("Do: ", err)
        return
    }

    // defer closing
    defer resp.Body.Close()

    // put response (which should be json) into a byte array for use by json.Unmarshal
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    b := buf.Bytes()

    // Simple check to see if we don't have json
    if b[0] != '{' {
        log.Fatal(fmt.Sprintf("ERROR: No json received -> %v", buf))
    }

    // Populate struct with json data
    var serverlist SteamServerList

    if err := json.Unmarshal(b, &serverlist); err != nil {
        log.Fatal("Unmarshal: ", err)
    }

    // Iterate over servers and print IP and query port
    for _, server := range serverlist.Response.Servers {
        tokens := strings.Split(server.Addr, ":")
        fmt.Printf("%s %s\n", tokens[0], tokens[1])
    }
}
