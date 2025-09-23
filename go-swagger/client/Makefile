generate:
	rm -rf ./generated
	mkdir generated
	swagger generate client -f ../swagger.yaml -t ./generated
	go mod tidy