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
    "sort"
    "regexp"
)

// nested struct to hold json data
type SteamServerList struct {
    Response `json:"response"`
}
type Response struct {
    Servers []Server `json:"servers"`
}
type Server struct {
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
}

// implement sorting interface
type ServerList []Server
func (a ServerList) Len() int { return len(a) }
func (a ServerList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ServerList) Less(i, j int) bool { return a[i].Name < a[j].Name; }

func main() {
    // Get command line arguments
    KeyPtr := flag.String("key", "", "Steam API key (**REQUIRED**)")
    LimitPtr := flag.Int("limit", 5000, "Limit search results")
    FilterPtr := flag.String("filter", "", "filter string")
    PlayersPtr := flag.Bool("players", false, "show player info")
    DebugPtr := flag.Bool("debug", false, "show debug output")
    DisplayPtr := flag.Bool("display", false, "Display full server info table")
    KickersPtr := flag.Bool("kickers", false, "Display server info with kick in name")
    flag.Parse()

    // Steam API key is required
    if *KeyPtr == "" {
        flag.PrintDefaults()
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

    // sort servers by name
    var serverList ServerList
    serverList = serverlist.Response.Servers
    sort.Sort(serverList)

    playerCount := 0

    // Iterate over servers and print IP and query port
    for _, server := range serverList {
        if *DebugPtr {
            fmt.Printf("%+v\n", server)
        }

        if *DisplayPtr {
            r, _ := regexp.Compile("[0-9]{1,2}:[0-9]{1,2}")
            Time := r.Find([]byte(server.Gametype))
            Perspective := "3PP"
            if strings.Contains(server.Gametype, "no3rd") {
                Perspective = "1PP"
            }
            serverName := server.Name
            if len(serverName) > 52 {
                serverName = serverName[:52]
            }
            if server.Appid == 221100 {
                fmt.Printf("%-52s %2d/%2d %s %s %s\n", serverName, server.Players, server.MaxPlayers, Time, Perspective, server.Version)
            } else {
                fmt.Printf("%-52s %2d/%2d %s %s\n", serverName, server.Players, server.MaxPlayers, Time, server.Version)
            }
            playerCount += server.Players
        } else {
            tokens := strings.Split(server.Addr, ":")
            if *PlayersPtr {
                fmt.Printf("%s %s %d %d\n", tokens[0], tokens[1], server.Players, server.MaxPlayers)
            } else if *KickersPtr {
                fmt.Printf("%-15s\t%d\t%s\n", tokens[0], server.Gameport, server.Name)
            } else {
                fmt.Printf("%s %s\n", tokens[0], tokens[1])
            }
        }
    }

    if *DisplayPtr {
        fmt.Printf("\n%d players on %d servers\n", playerCount, len(serverList))
    }
}
