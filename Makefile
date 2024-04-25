OUT=gogr

.PHONY: lint
lint:
	@golint ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: build
build: vet
	@go build -o artifacts/$(OUT) .

.PHONY: install
install: build
	@go install ./...

.PHONY: test
test: vet
	go test ./...

.PHONY: coverage
coverage: vet
	@./tools/coverage.sh

.PHONY: coverage-html
coverage-html: vet
	@./tools/coverage.sh html
