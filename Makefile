PKGS := $(shell go list ./...)

test:
	go test -v $(PKGS)

lint:
	golangci-lint run -v

vet:
	go vet $(PKGS)

sec:
	gosec ./...

fmt-check:
	goimports -l *.go **/*.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	goimports -w *.go **/*.go

lint-fix:
	golangci-lint run -v --fix

fix:
	$(MAKE) fmt
	$(MAKE) lint-fix
