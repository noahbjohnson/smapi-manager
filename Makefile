build: clean
	wails build

run: clean
	go run main.go

clean:
	rm -rf build/
	rm -rf frontend/build/

lint:
	cd frontend && npm run lint
	go fmt

compile:
	wails build -p
	wails build -x linux/amd64 -p
	wails build -x linux/arm-7 -p
	wails build -x windows/amd64 -p