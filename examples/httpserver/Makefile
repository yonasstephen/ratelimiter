BINARY_NAME?=httpserver

build:
	go build -o out/$(BINARY_NAME)

clean:
	rm -rf ./out

start: clean build
	chmod +x ./out/$(BINARY_NAME)
	./out/$(BINARY_NAME)