GO_FLAGS   = 
EXECUTABLE = tasiadb

main.go:
	go build $(GO_FLAGS) -o bin/$(EXECUTABLE) src/main.go

run:
	go run src/main.go
