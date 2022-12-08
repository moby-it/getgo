.PHONY: build
build: 
	go build -o bin/ main.go

.PHONY: test
test: 
	go test ./...

.PHONY: run
run: 
	go run main.go