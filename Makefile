
tidy:
	go mod tidy

fmt: tidy
	goimports -w .

lint: fmt
	golangci-lint run .

.PHONY: tidy fmt lint 
