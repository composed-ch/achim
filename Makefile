.PHONY: all format clean

all: achim

achim: cmd/achim/main.go
	go build -o achim cmd/achim/main.go

format:
	goimports -w .

clean:
	rm -f achim
