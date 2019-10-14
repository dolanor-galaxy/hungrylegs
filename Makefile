.PHONY: build clean

run:
	CGO_ENABLED=1 go run src/main.go

run.docker:
	docker run --rm -it -p 8000:8000 therohans/hungrylegs

test:
	cd src; go test ./...

clean:
	rm -rf build

build.docker: clean
	docker build -t therohans/hungrylegs .

build:
	mkdir build
	CGO_ENABLED=1 go build -o hungrylegs src/main.go 
	mv ./hungrylegs build/
	cp ./config.prod.json build/config.json
	cp -R ./store/ build/
	cp -R ./import/ build/
	cp -R ./migrations/ build/
