.PHONY: build clean

run:
	CGO_ENABLED=1 go run src/main.go

test:
	cd src; go test ./...

clean:
	rm -rf build

build:
	mkdir build
	CGO_ENABLED=1 go build -o hungrylegs src/main.go 
	mv hungrylegs build
	cp config.json build
	cp -R store build
	mkdir build/import
	cp -R migrations build