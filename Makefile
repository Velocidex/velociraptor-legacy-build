all:
	go run -v ./make.go build

clean:
	rm -rf ./build/
