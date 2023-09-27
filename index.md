# rp
Simple Reverse Proxy in Go

## Installation
Download pre-build version from [here](https://github.com/mchrisson/rp/releases)

OR

```
git clone https://github.com/mchrisson/rp.git
```

## Build
```
go build
```

## Usage
```
./rp -d 'https://remote.server' -p '7869'
```

## Notes
Also obeys proxy environment variables _(HTTP_PROXY, HTTPS_PROXY, NO_PROXY)_
