# iptracking
Get OS TCP Connection State

iptracking旨在提供一个实时追踪站外站内IP地址链路的简化工具。

## Info
- iptracking V0.1.X目前只实现Linux下地址端口解析呈现功能，不建议在任何线上环境使用。如果发现工具异常问题，请尽快留言，我们将尽快修复。
- iptracking 目前并没有创建官网。

## Features

- Support for Linux/mac/windows
- GeoIP rule support
- Configuration
- IPtracking Dashboard

## Install

You can download from Release page

## Build
- Make sure have python3 and golang installed in your computer.

- Install Golang
  ```
  brew install golang

  or download from https://golang.org
  ```

- Download deps

- Build and run.
  ```
  Linux Build And run
  GOOS=linux   GOARCH=amd64 go build -o iptracking-amd64-linux  iptracking.go
  ./iptracking-amd64-linux

  Darwin Build And run
  GOOS=darwin   GOARCH=amd64 go build -o iptracking-amd64-darwin  iptracking.go
  ./iptracking-amd64-darwin

  Windows Build And run
  GOOS=windows   GOARCH=amd64 go build -o iptracking-amd64-windows  iptracking.go
  iptracking-amd64-windows.exe

  ```


### FAQ
- Q: N/A?  
  A: N/A
