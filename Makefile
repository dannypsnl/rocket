docs:
	hugo --source ./rocket-doc --destination docs
	mv ./rocket-doc/docs .

GO_MOD=GO111MODULE=on

.PHONY: deps test coverage
deps:
	@$(GO_MOD) go get ./...
test: deps
	@$(GO_MOD) go test ./... -count 1 -cover -failfast
coverage: deps
	@$(GO_MOD) go test -coverprofile=coverage.txt ./...