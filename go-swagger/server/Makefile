generate:
	rm -rf ./generated
	mkdir generated
	swagger generate server -f ../swagger.yaml -t ./generated --exclude-main
	go mod tidy