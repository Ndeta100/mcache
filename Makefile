build:
	@go build -o mcache

run:build
	@./mcache

test:
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ./mcache

docker:
	echo "building docker file"
