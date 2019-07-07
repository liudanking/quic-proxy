# Quic Proxy

A http/https proxy using [QUIC](https://www.chromium.org/quic) as transport layer.

## Why use QUIC as transport layer instead of TCP?

* Almost 0-RTT for connection establishment
* Multiplexing
* Improved congestion control
* FEC
* Connection migration

Implementation detail: [A http proxy based on QUIC in 100 lines](https://liudanking.com/beautiful-life/100%E8%A1%8C%E4%BB%A3%E7%A0%81%E5%AE%9E%E7%8E%B0%E5%9F%BA%E4%BA%8E-quic-%E7%9A%84-http-%E4%BB%A3%E7%90%86/).

## Architecture 

![](https://ws1.sinaimg.cn/large/44cd29dagy1fpn4yaf2p8j20nd079aae.jpg)

## Installation & Usage

**Note**: require go version >= 1.9

### Install `qpserver` on your remote server

`go get -u github.com/liudanking/quic-proxy/qpserver`

### Start `qpserver`:

`qpserver -v -l :3443 -cert YOUR_CERT_FILA_PATH -key YOUR_KEY_FILE_PATH -auth username:password`

### Install `qpclient` on your local machine

`go get -u github.com/liudanking/quic-proxy/qpclient`

### Start `qpclient`:

`qpclient -v -k -proxy http://YOUR_REMOTE_SERVER:3443 -l 127.0.0.1:18080 -auth username:password`

### Set proxy for your application on your local machine

Let's take Chrome with SwitchyOmega for example:

![](https://ws1.sinaimg.cn/large/44cd29dagy1fpn5c4jng6j21eq0fw40j.jpg)

Enjoy!

## TODO

* Using custom congestion control

## Join Wechat Group

Add the Wechat robot to join group:

<img src="https://raw.githubusercontent.com/liudanking/quic-proxy/master/wx-bot.jpg" width="320px"/>

