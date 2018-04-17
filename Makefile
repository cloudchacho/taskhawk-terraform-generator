.PHONY: test

test_setup:
	./scripts/test-setup.sh

test: clean test_setup
	./scripts/run-tests.sh

build:
	go-bindata -prefix "templates/" templates/

	env GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/taskhawk-terraform-generator .
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/taskhawk-terraform-generator .

clean:
	rm -rf bin
