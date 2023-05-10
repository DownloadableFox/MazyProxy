name = mazy-proxy

build:
	@go build -o bin/$(name) src/**.go

run: build
	@./bin/$(name)