generate:
	rm -rf ./generated go.mod go.sum
	go mod init client || true
	mkdir generated
	swagger generate client -f ../swagger.yaml -t ./generated
	go mod tidy