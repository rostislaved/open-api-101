generate:
	rm -rf ./generated go.mod go.sum
	go mod init client || true
	mkdir generated
	ogen --target generated --clean ../openapi.yaml
	go mod tidy
