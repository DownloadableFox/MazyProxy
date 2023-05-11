name = mazy-proxy.exe

build:
	@go build -o bin/$(name) src/**.go

run: build
	@./bin/$(name)