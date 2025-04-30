#!/bin/sh

rm ./HipparchiaGoBuilder

GO="go"
DT=$(date "+%Y-%m-%d@%H:%M:%S")
GC=$(git rev-list -1 HEAD | cut -c-8)

PG="default.pgo"
PGF="${PG}"

LDF="-s -w -X main.BuildDate=${DT}"

${GO} build -pgo=${PGF} -ldflags "${LDF}"
