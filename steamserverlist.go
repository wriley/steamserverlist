package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// nested struct to hold json data
type steamJSON struct {
	response `json:"response"`
}
type response struct {
	Servers []server `json:"servers"`
}
type server struct {
	Addr       string `json:"addr"`
	Gameport   int    `json:"gameport"`
	Steamid    string `json:"steamid"`
	Name       string `json:"name"`
	Appid      int    `json:"appid"`
	Gamedir    string `json:"gamedir"`
	Version    string `json:"version"`
	Product    string `json:"product"`
	Region     int    `json:"region"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Bots       int    `json:"bots"`
	Map        string `json:"map"`
	Secure     bool   `json:"secure"`
	Dedicated  bool   `json:"dedicated"`
	Os         string `json:"os"`
	Gametype   string `json:"gametype"`
}

// implement sorting interface
type serverList []server

func (a serverList) Len() int           { return len(a) }
func (a serverList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a serverList) Less(i, j int) bool { return a[i].Name < a[j].Name }

func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

func addCommaSeperator(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func main() {
	// Get command line arguments
	KeyPtr := flag.String("key", "", "Steam API key (**REQUIRED**)")
	LimitPtr := flag.Int("limit", 10000, "Limit search results")
	FilterPtr := flag.String("filter", "", "filter string")
	PlayersPtr := flag.Bool("players", false, "show player info")
	DebugPtr := flag.Bool("debug", false, "show debug output")
	DisplayPtr := flag.Bool("display", false, "Display full server info table")
	Display2Ptr := flag.Bool("display2", false, "Display full server info table and IP/Port")
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
	timeout := time.Duration(15 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
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
	var serverJSON steamJSON

	if err := json.Unmarshal(b, &serverJSON); err != nil {
		log.Fatal("Unmarshal: ", err)
	}

	// sort servers by name
	var allServers serverList
	allServers = serverJSON.response.Servers
	sort.Sort(allServers)

	playerCount := 0
	queueCount := 0

	// Iterate over servers and print IP and query port
	for _, server := range allServers {
		if *DebugPtr {
			fmt.Printf("%+v\n", server)
		}

		if *DisplayPtr || *Display2Ptr {
			r, _ := regexp.Compile("[0-9]{1,2}:[0-9]{1,2}")
			Time := r.Find([]byte(server.Gametype))
			TimeMultiplier := 1.0
			TimeMultiplierNight := 1.0

			r, _ = regexp.Compile("etm([0-9]{1,3}.[0-9]{1,6})")
			TimeMultiMatches := r.FindStringSubmatch(server.Gametype)
			if len(TimeMultiMatches) > 0 {
				i, err := strconv.ParseFloat(TimeMultiMatches[1], 64)
				if err == nil {
					TimeMultiplier = i
				}
			}

			r, _ = regexp.Compile("entm([0-9]{1,3}.[0-9]{1,6})")
			TimeMultiMatches = r.FindStringSubmatch(server.Gametype)
			if len(TimeMultiMatches) > 0 {
				i, err := strconv.ParseFloat(TimeMultiMatches[1], 64)
				if err == nil {
					TimeMultiplierNight = i
				}
			}

			TimeMultiplierString := fmt.Sprintf("%.1fx/%.1fx", TimeMultiplier, TimeMultiplierNight)
			TimeMultiplierString = strings.Replace(TimeMultiplierString, ".0x", "x", -1)

			QueueSize := int64(0)
			r, _ = regexp.Compile("lqs([0-9]+)")
			QueueSizeMatches := r.FindStringSubmatch(server.Gametype)
			if len(QueueSizeMatches) > 0 {
				i, err := strconv.ParseInt(QueueSizeMatches[1], 10, 64)
				if err == nil {
					QueueSize = i
				}
			}
			queueCount += int(QueueSize)

			Perspective := "3PP"
			if strings.Contains(server.Gametype, "no3rd") {
				Perspective = "1PP"
			}
			serverName := stripCtlAndExtFromBytes(server.Name)
			if server.Appid == 221100 || server.Appid == 1024020 {
				if len(serverName) > 50 {
					serverName = serverName[:50]
				}
				if *Display2Ptr {
					tokens := strings.Split(server.Addr, ":")
					fmt.Printf("%-50s %3d/%-3d %s %s %s %s %s %s\n", serverName, server.Players, server.MaxPlayers, Time, TimeMultiplierString, Perspective, server.Version, tokens[0], tokens[1])
				} else {
					fmt.Printf("%-50s %3d/%-3d %s %s %s %s\n", serverName, server.Players, server.MaxPlayers, Time, TimeMultiplierString, Perspective, server.Version)
				}
			} else {
				if len(serverName) > 52 {
					serverName = serverName[:52]
				}
				fmt.Printf("%-56s %2d/%2d %s %s\n", serverName, server.Players, server.MaxPlayers, Time, server.Version)
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

	if *DisplayPtr || *Display2Ptr {
		fmt.Printf("\n%s players on %s servers and %s in queue\n", addCommaSeperator(int64(playerCount)), addCommaSeperator(int64(len(allServers))), addCommaSeperator(int64(queueCount)))
	}
}
