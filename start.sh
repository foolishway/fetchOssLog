#!/bin/bash
rm ./fetchosslog
go build -o fetchosslog *.go
chmod +x ./fetchosslog
nohup ./fetchosslog > log.txt 2>&1 &
