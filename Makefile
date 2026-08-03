.PHONY: run format clean check test cover release clean

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

release: cmd/achim/achim.go
	mkdir -p release/linux release/macos release/windows

	GOARCH=amd64 GOOS=linux go build -o release/linux/achim cmd/achim/achim.go
	GOARCH=arm64 GOOS=darwin go build -o release/macos/achim cmd/achim/achim.go
	GOARCH=amd64 GOOS=windows go build -o release/windows/achim.exe cmd/achim/achim.go

	zip release/achim-linux-amd64.zip release/linux/achim
	zip release/achim-macos-arm64.zip release/macos/achim
	zip release/achim-windows-amd64.zip release/windows/achim.exe

clean:
	rm -f cover.html cover.out achim achim.exe
	rm -rf release
