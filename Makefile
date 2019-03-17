run_server:
	go run main.go --server_port 40059 --client_port 40000 --client_port 40001

run_client:
	go run main.go --server_port 40000 --server_port 40001 --client_port 30000

build:
	go build -o bin/udp_proxy main.go

default: build
