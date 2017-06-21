version=v0.1.0
binary_tarball="cook-$(version).tar.gz"
binary_location="${HOME}/bin/cook"
bash_completion_dir ?= "${HOME}/.bash_completion.d/cook.bash-completion"

.PHONY: get-deps
get-deps:
	go get github.com/sirupsen/logrus

.PHONY: build
build: main.go search.go validate.go parser.go
	go build -ldflags "-X main.version=$(version)" -o cook main.go search.go validate.go parser.go

.PHONY: local
local: build
	install cook "$(binary_location)"

.PHONY: install-bash-completion
install-bash-completion:
	@install -v -d "$(bash_completion_dir)" && \
		install -m 0644 -v completion/cook.bash-completion "$(bash_completion_dir)/cook"

.PHONY: archive
archive: build
	tar -cvzf "$(binary_tarball)" cook
