all: test kakedit kakpipe

kakedit:
	go build -o kakedit cmd/kakedit/main.go

kakpipe:
	go build -o kakpipe cmd/kakedit/main.go

test:
	go vet ./...
	go test ./...

install:
	go install cmd/kakedit/main.go
	go install cmd/kakpipe/main.go
