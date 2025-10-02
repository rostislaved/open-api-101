generate:
	rm -rf ./generated go.mod go.sum
	go mod init server || true
	mkdir generated
	oapi-codegen -config cfg.yaml ../../openapi.yaml
	go mod tidy

mockery:
	rm -f ./handlers/mocks_test.go
	mockery
	go mod tidy