.PHONY: build
build: main.go search.go
	go build -o cook main.go search.go

.PHONY: install
install:
	install cook "${HOME}/bin/cook"
