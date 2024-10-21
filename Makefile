
tidy:
	go mod tidy

fmt: tidy
	goimports -w .
