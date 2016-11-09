# steamserverlist

A Steam server list program written in [golang](https://golang.org/) that uses the [Steam Web API](https://partner.steamgames.com/documentation/webapi)

## Example Usage

### DayZ servers with "hardcore" in the name
```
$ steamserverlist -key YOURAPIKEY -filter '\gamedir\dayz\name_match\*hardcore*'
103.13.101.247:27700 2500
64.94.95.98:27500 2300
64.94.95.154:27700 2500
199.60.101.202:27600 2400
66.55.158.194:27017 2302
192.3.53.26:27016 2302
108.174.57.251:27116 2402
108.61.112.232:27017 2302
81.19.216.161:27600 2400
109.70.149.150:27700 2500
81.19.216.158:27520 2320
37.187.72.43:27516 2802
185.62.204.42:27316 2602
```
