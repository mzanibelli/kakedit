SRC = $(wildcard *.go)
SRC += $(wildcard internal/kakoune/*.go)
SRC += $(wildcard internal/listener/*.go)

all: test kakedit kakpipe kakwrap

kakedit: $(SRC) cmd/kakedit/main.go
	go build -o kakedit cmd/kakedit/main.go

kakpipe: $(SRC) cmd/kakpipe/main.go
	go build -o kakpipe cmd/kakpipe/main.go

kakwrap: $(SRC) cmd/kakwrap/main.go
	go build -o kakwrap cmd/kakwrap/main.go

test:
	go vet ./...
	go test ./...

install:
	go install ./cmd/kakedit
	go install ./cmd/kakpipe
	go install ./cmd/kakwrap

clean:
	rm -f kakedit kakpipe kakwrap
