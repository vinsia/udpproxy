run_server:
	go run main.go -s 127.0.0.1:40059 -c 127.0.0.1:40000 -c 127.0.0.1:40001

run_client:
	go run main.go -s 127.0.0.1:40000 -s 127.0.0.1:40001 -c 127.0.0.1:30000 -l

build:
	go build -o bin/udp_proxy main.go

default: build
