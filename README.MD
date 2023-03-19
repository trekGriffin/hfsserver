[![Release Binary](https://github.com/trekGriffin/hfsserver/actions/workflows/go2.yml/badge.svg)](https://github.com/trekGriffin/hfsserver/actions/workflows/go2.yml)

# simple http file system  with curl
usage: hfsserver port directory

## get a file
1. open url in the browser: http://ip:port/

2. curl curl -O http://ip:port/1.txt


## post a file
curl -T  1.txt http://ip:port/upload/1.txt

## delete a file
curl http://ip:port/delete/1.txt