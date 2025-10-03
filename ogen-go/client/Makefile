generate:
	rm -rf ./generated go.mod go.sum
	go mod init client || true
	mkdir generated
	ogen --config ./ogen.yaml --target generated --clean ../openapi.yaml
	go mod tidy
