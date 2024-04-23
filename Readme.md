# ltp API


## How to run 
You need create config.yaml like:
```yaml
server:
  host: "0.0.0.0"
  port: 84
  readTimeout: 30s
  writeTimeout: 30s
  swaggerAddress: "localhost:84"

redis:
  host: "193.57.209.218"
  port: 6379
  password: ""

kraken:
  url: https://api.kraken.com
  withLogs: true

pairs:
  - BTC/USD
  - BTC/CHF
  - BTC/EUR
```
After create docker-compose.yml
```yaml
version: "3.9"

services:
  ltp:
    container_name: ltp-service
    build:
      context: .
    ports:
      - "84:84"
    volumes:
      - ./config.yaml:/etc/config.yaml
    restart: unless-stopped
```
After run 
```shell
docker-compose up -d
```
After you get log like
```shell
[+] Running 2/2
 ✔ Network ltp-api_default  Created                                                                                                                                                                                            0.0s 
 ✔ Container ltp-service    Started      
```


After you can open http://localhost:84/docs/index.html#/LTPHandler/get_ltp and check the API



Project uses modular monolith architecture.

## Project Structure
```
├───cmd
├───docs
├───internal
│   ├───bootstrap
│   ├───config
│   ├───constant
│   ├───database
│   │   ├───entities
│   │   └───repositories
│   ├───dto
│   ├───services
│   ├───transport
│   │   └───http
│   │       ├───handlers
│   │       ├───middlewares
│   │       ├───requests
│   │       ├───response
│   │       └───server
│   ├───utils
│   └───validator
├───migrations
├───pkg
│   ├───database
│   ├───kanel
│   ├───map_structure
│   ├───redis
│   └───smtp
└───templates
```
