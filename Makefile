all: test kakedit kakpipe kakwrap

kakedit:
	go build -o kakedit cmd/kakedit/main.go

kakpipe:
	go build -o kakpipe cmd/kakedit/main.go

kakwrap:
	go build -o kakwrap cmd/kakedit/main.go

test:
	go vet ./...
	go test ./...

install:
	go install ./cmd/kakedit
	go install ./cmd/kakpipe
	go install ./cmd/kakwrap
