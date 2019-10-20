.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)
os = macOS10.13

# List all targets in thie file
list:
	@echo ""
	@echo "🏊🏻‍🚴‍🏃‍ HungryLegs 🏊🏻‍🚴‍🏃‍"
	@echo ""
	@grep -B 1 '^[^#[:space:]].*:' Makefile
	@echo ""

# Run all go unit tests
test:
	go test ./...

# Runs the graphql endpoint server (localhost)
run.server:
	CGO_ENABLED=1 go run cmd/server/server/server.go

# Run the local cli (will use config.json)
run.cli:
	CGO_ENABLED=1 go run cmd/cli/main.go

# Run the plan application (.csv to .ics application)
run.plan:
	go run cmd/plan/main.go ./testdata/example.csv ./testdata/plan.ics

# Build the plan application (.csv to .ics application)
build.plan:
	mkdir -p build
	go build -o build/plan cmd/plan/main.go

# Remove build artifacts
clean:
	rm -rf build

# Builds a version within Docker (Linux)
build.docker: clean
	docker build -t therohans/hungrylegs .

# Builds a local OS version
build.cli: clean
	mkdir -p build
	CGO_ENABLED=1 go build -o hungrylegs cmd/cli/main.go 
	mv ./hungrylegs build/
	cp ./config.json build/config.json
	mkdir -p build/store/athletes
	mkdir -p build/import/
	cp -R ./migrations/ build/migrations/

# Builds the graphql endpoint server
build.server: clean
	mkdir -p build
	CGO_ENABLED=1 go build -o hungrylegs cmd/server/server/server.go 
	mv ./hungrylegs build/
	cp ./config.prod.json build/config.json
	mkdir -p build/store/athletes
	cp -R ./import/ build/
	cp -R ./migrations/ build/

# Run the Dockerfile (only need when working on the dockerfile itself)
run.docker:
	docker run --rm -it -p 8080:8080 therohans/hungrylegs

# Build all the tools for local playing aroundc
build.all.cli: clean build.cli build.plan

# Create a zip file for the command line tools
distrib.cli: build.all.cli
	mv build hungrylegs
	zip hungrylegs-$(os)-$(hash).zip -r ./hungrylegs/*
	rm -rf hungrylegs