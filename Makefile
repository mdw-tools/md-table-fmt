.PHONY: test build install clean

test:
	go mod tidy
	go fmt ./...
	go test -race -cover ./...

build:
	go build -o md-table-fmt .

install:
	go install .

clean:
	rm -f md-table-fmt
