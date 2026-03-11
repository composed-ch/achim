.PHONY: all run format clean check test cover clean

all: achim achim.exe

achim: cmd/achim/achim.go
	go build -o achim $^

achim.exe: cmd/achim/main.go
	GOOS=windows go build -o achim.exe cmd/achim/main.go

run: cmd/achim/achim.go
	go run $^

format:
	goimports -w .

check:
	go vet ./...

test:
	go test -cover ./...

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	rm -f cover.out

clean:
	rm -f cover.html cover.out achim achim.exe
