.PHONY: check tidy deps \
	lint lint-md lint-go \
	lint-fix lint-md-fix

check:
	go test ./...

tidy:
	go mod tidy

deps:
	go install github.com/mgechev/revive@latest
	go install golang.org/x/tools/cmd/goimports@latest
	npm install

lint: lint-md lint-go
lint-fix: lint-fix-md

lint-md:
	npx remark . .github

lint-fix-md:
	npx remark . .github -o

lint-go:
	revive -formatter stylish -config revive.toml ./...
