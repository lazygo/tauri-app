#!/bin/sh
cd "$(dirname "$0")"
echo "[`date`]: $1/lib$2.a" >> build.log
go build -buildmode=c-archive -trimpath -o $1/lib$2.a module.go >> build.log 2>&1