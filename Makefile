
install_dependencies:
	go get ./...

dev:
	go run src/main/*.go

build:
	go build -o audigo_stream src/main/*.go