.PHONY: build clean test

# List all targets in thie file
list:
	@echo ""
	@echo "🏊🏻‍🚴‍🏃‍ HungryLegs 🏊🏻‍🚴‍🏃‍"
	@echo ""
	@grep -B 1 '^[^#[:space:]\.].*:' Makefile
	@echo ""

# Run all go unit tests
test:
	go test ./...

# Run the local cli (will use config.json)
run.cli:
	CGO_ENABLED=1 go run cmd/cli/main.go

# Remove build artifacts
clean:
	rm -rf build

# Builds a version within Docker (Linux)
build.docker: clean
	docker build -t therohans/hungrylegs .

# Builds a local OS version
build: clean
	mkdir build
	CGO_ENABLED=1 go build -o hungrylegs cmd/cli/main.go 
	mv ./hungrylegs build/
	cp ./config.prod.json build/config.json
	cp -R ./store/ build/
	cp -R ./import/ build/
	cp -R ./migrations/ build/

# Run the Dockerfile (only need when working on the dockerfile itself)
run.docker:
	docker run --rm -it -p 8000:8000 therohans/hungrylegs
