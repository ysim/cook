version=v0.2.0
binary_tarball="cook-$(version).tar.gz"
binary_location="${HOME}/bin"
bash_completion_dir ?= "${HOME}/.bash_completion.d"

.PHONY: build
build: main.go search.go validate.go parser.go create.go list.go util.go display.go
	go build -ldflags "-X main.version=$(version)" -o cook $^
	go build -ldflags "-X main.version=$(version)" -o concoct $^

.PHONY: local
local: build
	install cook "$(binary_location)/cook"
	install concoct "$(binary_location)/concoct"

.PHONY: install-bash-completion
install-bash-completion:
	@install -v completion/bash_completion "${HOME}/.bash_completion"
	@install -v -d "$(bash_completion_dir)" && \
		install -m 0644 -v completion/cook.bash-completion "$(bash_completion_dir)/cook.bash-completion"

.PHONY: release-darwin
release-darwin:
	$(MAKE) build GOOS=darwin GOARCH=amd64
	@mkdir -p release/"$(version)"
	tar -cvzf release/"$(version)"/"cook-darwin-$(version).tar.gz" cook
	shasum -a 256 release/"$(version)"/"cook-darwin-$(version).tar.gz" | awk '{print $$1}' > release/"$(version)"/cook-darwin-$(version).sha256

.PHONY: release-linux
release-linux:
	$(MAKE) build GOOS=linux GOARCH=amd64
	@mkdir -p release/"$(version)"
	tar -cvzf release/"$(version)"/"cook-linux-$(version).tar.gz" cook
	shasum -a 256 release/"$(version)"/"cook-linux-$(version).tar.gz" | awk '{print $$1}' > release/"$(version)"/cook-linux-$(version).sha256
