run_server:
	go run main.go -s 40059 -c 40000 -c 40001 -l

run_client:
	go run main.go -s 40000 -s 40001 -c 30000 -l

build:
	go build -o bin/udp_proxy main.go

default: build
