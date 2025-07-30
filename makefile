.PHONY: run build test clean fmt debug-run

run:
	go run main.go

debug-run:
	go run main.go -debug=true

build:
	go build -o glog main.go

test:
	go test ./...

clean:
	go clean
	rm -f glog

fmt:
	go fmt ./...