run:
	@go run main.go

build:
	@echo "Building for development"
	@go build -o ./bin/build/test-url-tls main.go
	@echo "Done"

clean:
	@echo "Cleaning..."
	@rm ./bin/build/* || echo "Empty directory"
	@rm ./bin/release/* || echo "Empty directory"
	@echo "Done"

release:
	@rm ./bin/release/* || echo "New build! Yay!"
	@echo "Building for release: linux/amd64"
	@GOOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o ./bin/release/test-url-tls main.go && gzip -c ./bin/release/test-url-tls > ./bin/release/test-url-tls-linux-amd64.gz
	@echo "Building for release: windows/amd64"
	@GOOARCH=amd64 GOOS=windows go build -ldflags="-w -s" -o ./bin/release/test-url-tls.exe main.go && gzip -c ./bin/release/test-url-tls.exe > ./bin/release/test-url-tls-windows-amd64.gz
	@rm ./bin/release/test-url-tls
	@rm ./bin/release/test-url-tls.exe
	@echo "Done"
