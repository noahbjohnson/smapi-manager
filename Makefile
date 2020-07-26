build: clean bundle_frontend
	go build -o bin/main main.go

run: clean bundle_frontend
	go run main.go

clean:
	rm -rf bin/
	rm -rf statik/

install:
	go get github.com/rakyll/statik

install_frontend:
	cd frontend && npm install

lint_frontend:
	cd frontend && npm run lint

build_frontend: install_frontend
	cd frontend && npm run build

bundle_frontend: build_frontend install
	statik -src=frontend/build/


compile: bundle_frontend
	# 32-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
    # MacOS
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main.go
    # Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
    # Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
    # 64-Bit
    # FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go
    # MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
    # Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
    # Winodws
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go