#!/bin/bash

go mod tidy
go run ./cmd/server $@
