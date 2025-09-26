generate:
	rm -rf ./generated go.mod go.sum
	go mod init server || true
	mkdir generated
	swagger generate server -f ../swagger.yaml -t ./generated --exclude-main
	go mod tidy

mockery:
	rm -f ./handlers/mocks_test.go
	mockery
	go mod tidy
