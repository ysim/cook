.PHONY: build
build: main.go search.go validate.go
	go build -o cook main.go search.go validate.go

.PHONY: install
install: build
	install cook "${HOME}/bin/cook"
