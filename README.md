# steamserverlist

A Steam server list program written in [golang](https://golang.org/) that uses the [Steam Web API](https://partner.steamgames.com/documentation/webapi)

## Example Usage

### DayZ servers with "hardcore" in the name
```
$ steamserverlist -key YOURAPIKEY -filter '\gamedir\dayz\name_match\*hardcore*'
109.70.149.150 27700
173.199.67.66 27017
185.62.204.42 27316
103.13.101.247 27700
64.94.95.98 27500
108.174.57.251 27116
37.187.72.43 27516
81.19.216.158 27520
66.55.158.194 27017
108.61.112.232 27017
81.19.216.161 27600
64.94.95.154 27700
```
