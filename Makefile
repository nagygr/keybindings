.phony: all chk

chk:
	go fmt -mod=mod ./cmd/...
	go vet ./cmd/...

all: chk
	go build -mod=mod -o build/keybindings ./cmd/keybindings
