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

# update
* 230407.1
fix to real vcs
get file name from the local file other than from the url.
show the precise speed of uploading files.

* 1.0.2
bind the ip address for the software.(for security)
add the vcs to this software
auto generate linux binary
* 1.0.1
* show uploading progress bar for curl and server
* add changelog and README

