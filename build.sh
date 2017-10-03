#!/usr/bin/env bash

   go generate \
&& go fmt templates.go >/dev/null \
&& go build -o createESP32App
