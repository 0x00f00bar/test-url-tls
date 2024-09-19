run:
	go run main.go

build:
	go build -o ./bin/build/test-url-tls main.go

release:
	go build -ldflags="-w -s" -o ./bin/release/test-url-tls main.go
