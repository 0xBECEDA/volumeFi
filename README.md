# Flights tracker 

### Story
There are over 100,000 flights a day, with millions of people and cargo being transferred around the world. With so many people and different carrier/agency groups, it can be hard to track where a person might be. In order to determine the flight path of a person, we must sort through all of their flight records.

### Goal
To create a simple microservice API that can help us understand and track how a particular person's flight path may be queried. The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.

Required JSON structure:
```
[["SFO", "EWR"]]                                                 => ["SFO", "EWR"]
[["ATL", "EWR"], ["SFO", "ATL"]]                                 => ["SFO", "EWR"]
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]
```

### Implementation idea explanation

We have corner cases:

1. Flights can be a round trip (Paris -> Milan, Milan->Paris). Flights pairs are unordered, so you can't find 
   start and end point. 

2. Flights can have cycles, but the start and end point are not round trip. For example path Berlin -> Milan, Milan->Madrid, Madrid-Milan, Milan-Paris 
   has a cycle, but we can find start (Berlin) and end (Paris) flight.

Also, we assume that:

1. Departure and destination airports can't be same (Milan-Milan) or empty.

2. Flights must be connected as a chain without gap. For example, flights chain "Berlin->Milan, Milan->Paris" is valid, 
   but "Berlin->Milan, Paris->Madrid" isn't. It means, that exists at least 1 path to visit all given airports. 

## How to run

1. For tests running, staying in project's root, run 
```make test```

2. To run app, staying in project's, run
```make build && make run```

It will build docker image and run app in docker container. 
Requirements: you should have installed docker 20+ version and docker compose 2+ version.

### How to interact with app
#### Successful examples
Easy trip:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["SFO", "ATL"]
       ]
}'
```
response:
```
 200 {"start-end-flights":["SFO","ATL"]}
```

Trip without cycles:
sfo -> atl -> gso -> ind -> ewr
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["IND", "EWR"],
        ["SFO", "ATL"],
        ["GSO", "IND"],
        ["ATL", "GSO"]
    ]
}'
```
response:
```
 200 {"start-end-flights":["SFO","EWR"]}
```

Trip with cycle in the end, but start and end point aren't round trip:
sfo -> atl -> gso -> ind -> ewr -> ind
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["IND", "EWR"],
        ["EWR", "IND"],
        ["SFO", "ATL"],
        ["GSO", "IND"],
        ["ATL", "GSO"]
    ]
}'
```
response:
```
 200 {"start-end-flights":["SFO","IND"]}
```

Trip with cycle in the beginning and in the end, but start and end point aren't round trip:
atl-> sfo-> atl-> gso -> ind -> ewr -> ind
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["IND", "EWR"],
        ["EWR", "IND"],
        ["ATL", "SFO"],
        ["SFO", "ATL"],       
        ["GSO", "IND"],
        ["ATL", "GSO"]
    ]
}'
```
response:
```
200 {"start-end-flights":["ATL","IND"]}
```

Trip with multiple cycles but start and end point aren't round trip:
ind-> gso-> ewd-> gso -> sfo -> alt -> sfo -> la
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["IND", "GSO"], 
        ["GSO", "EWR"],
        ["ATL", "SFO"], 
        ["EWR", "GSO"],
        ["SFO", "ATL"],       
        ["GSO", "SFO"], 
        ["SFO", "LA"]         
    ]
}'
```
response:
```
200 {"start-end-flights":["IND", "LA"]}
```

#### Fail examples
Same departure and destination airport:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["ATL", "ATL"]
       ]
}'
```
response:
```
 400 invalid airport codes
```

Empty destination airport:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["ATL", ""]
       ]
}'
```
response:
```
 400 invalid airport codes
```

Empty departure airport:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["", "ATL"]
       ]
}'
```
response:
```
 400 invalid airport codes
```

More than 2 airports: 
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["ALT", "BER", "SFO"]
       ]
}'
```
response:
```
 400 invalid airport codes
```

No airport codes:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": []
}'
```
response:
```
 400 request can't be empty
```

Number instead of strings array:
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [[1, "SFO"]]
}'
```
response:
```
400 invalid body: only strings are acceptable
```

Trip with cycle, start and end point are round trip:
sfo -> atl -> gso -> ind -> ewr -> ind -> sfo
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [
        ["IND", "SFO"],  
        ["IND", "EWR"], 
        ["EWR", "IND"], 
        ["SFO", "ATL"], 
        ["GSO", "IND"], 
        ["ATL", "GSO"] 
    ]
}'
```
response:
```
 400 route trip detected: unable find departure and arrival airport
```

Trip with cycle, but is not connected:
sfo -> atl -> sfo,  ind -> ewr
```
curl --location -i --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "flights": [        
        ["IND", "EWR"],         
        ["SFO", "ATL"],         
        ["ATL", "SFO"] 
    ]
}'
```
response:
```
 400 multiple paths detected: unable find departure and arrival airport
```

