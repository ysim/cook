.PHONY: build
build: main.go
	go build -o cook main.go

.PHONY: install
install:
	install cook "${HOME}/bin/cook"
