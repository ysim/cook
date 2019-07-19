version=v0.1.1
binary_tarball="cook-$(version).tar.gz"
binary_location="${HOME}/bin/cook"
bash_completion_dir ?= "${HOME}/.bash_completion.d"

.PHONY: build
build: main.go search.go validate.go parser.go
	go build -ldflags "-X main.version=$(version)" -o cook main.go search.go validate.go parser.go

.PHONY: local
local: build
	install cook "$(binary_location)"

.PHONY: install-bash-completion
install-bash-completion:
	@install -v completion/bash_completion "${HOME}/.bash_completion"
	@install -v -d "$(bash_completion_dir)" && \
		install -m 0644 -v completion/cook.bash-completion "$(bash_completion_dir)/cook.bash-completion"

.PHONY: archive
archive: build
	tar -cvzf "$(binary_tarball)" cook
	shasum -a 256 "$(binary_tarball)"
