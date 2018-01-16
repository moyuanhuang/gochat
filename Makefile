.PHONY: clean
default:
	go build -o gochat main.go
clean:
	go clean
