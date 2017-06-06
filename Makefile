.PHONY: get-deps
get-deps:
	go get github.com/sirupsen/logrus

.PHONY: build
build: main.go search.go validate.go parser.go
	go build -o cook main.go search.go validate.go parser.go

.PHONY: install
install: build
	install cook "${HOME}/bin/cook"
