.PHONY: all clean

all: achim

achim: cmd/achim/main.go
	go build -o achim cmd/achim/main.go

clean:
	rm -f achim
