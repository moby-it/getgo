BINPATH := \/home\/${USER}\/getgo\/bin\/getgo
HOMEPATH := \/home\/${USER}
CURRENT_USER := ${USER}
CURRENT_USER := ${GROUP}
.PHONY: clean
clean:
	rm -rf bin/
	go mod tidy

.PHONY: build
build: clean
	go build -o bin/getgo cmd/main.go

.PHONY:start
start: build
	./bin/main
.PHONY: test
test: 
	go test ./... -v

.PHONY: run
run: 
	go run cmd/main.go

.PHONY: install
install: build
	sudo ./install.sh