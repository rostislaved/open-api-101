generate:
	rm -rf ./generated go.mod go.sum
	go mod init client || true
	mkdir generated
	oapi-codegen -config cfg.yaml ../openapi.yaml
	go mod tidy
