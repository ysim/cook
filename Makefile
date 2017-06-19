version=v0.1.0
binary_location="${HOME}/bin/cook"

.PHONY: get-deps
get-deps:
	go get github.com/sirupsen/logrus

.PHONY: build
build: main.go search.go validate.go parser.go
	go build -ldflags "-X main.version=$(version)" -o cook main.go search.go validate.go parser.go

.PHONY: local
local: build
	install cook "$(binary_location)"
