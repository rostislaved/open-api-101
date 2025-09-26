generate:
	rm -rf ./generated go.mod go.sum
	go mod init server || true
	mkdir generated
	ogen --target generated --clean ../openapi.yaml
	go mod tidy

mockery:
	rm -f ./handlers/mocks_test.go
	mockery
	go mod tidy