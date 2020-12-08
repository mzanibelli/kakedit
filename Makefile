all: test kakedit kakpipe

kakedit:
	go build -o kakedit cmd/kakedit/main.go

kakpipe:
	go build -o kakpipe cmd/kakedit/main.go

test:
	go vet ./...
	go test ./...

install:
	go install ./cmd/kakedit
	go install ./cmd/kakpipe
