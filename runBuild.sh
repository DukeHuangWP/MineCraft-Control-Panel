#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./build/MineCraft-Control-Panel main.go

rm -rf ./website/downloads/ZIPs/*

#cp -rf ./deploy/* ./build/
#rm -rf ./build/configs/*
#cp -rf ./configs ./build
mkdir ./build/website
rm -rf ./build/configs/website/*
cp -rf ./website/* ./build/website
mkdir ./build/website/downloads/ZIPs
touch ./build/website/downloads/ZIPs/empty