.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)
os = macOS10.13

# echo 'export PATH="/usr/local/opt/sqlite/bin:$PATH"' >> ~/.bash_profile
# For compilers to find sqlite you may need to set:
LDFLAGS="-L/usr/local/opt/sqlite/lib"
CPPFLAGS="-I/usr/local/opt/sqlite/include"

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
	go run cmd/server/server/server.go

# Run the local cli
run.cli:
	CGO_ENABLED=1 \
	HL_DB_DRIVER=sqlite3 \
	HL_DB_CONNECTION=store/athletes/{athlete}.db \
	HL_DB_POST="PRAGMA synchronous = OFF;PRAGMA journal_mode = MEMORY;PRAGMA cache_size = -16000" \
	HL_BASE_IMPORT="/Users/robrohan/Dropbox/Documents/Fitness/CheetahAthlete/Professor Zoom/imports" \
	go run cmd/cli/main.go

# Run the plan application (.csv to .ics application)
run.plan:
	go run cmd/plan/main.go \
	--csv-path ./testdata/example.csv \
	--output ./testdata/plan.ics

# Build the plan application (.csv to .ics application)
build.plan:
	mkdir -p build
	go build -o build/plan -ldflags "-X main.build=${hash}" cmd/plan/main.go

# Remove build artifacts
clean:
	rm -rf build

# Builds a version within Docker (Linux)
build.docker: clean
	docker build -t robrohan/hungrylegs .

# Builds a local OS version
build.cli: clean
	mkdir -p build
	CGO_ENABLED=1 go build -o hungrylegs -ldflags "-X main.build=${hash}" cmd/cli/main.go 
	mv ./hungrylegs build/
	mkdir -p build/store/athletes
	mkdir -p build/import/
	cp -R ./migrations/ build/migrations/

# Builds the graphql endpoint server
build.server: clean
	mkdir -p build
	go build -o hungrylegs -ldflags "-X main.build=${hash}" cmd/server/server/server.go 
	mv ./hungrylegs build/
	mkdir -p build/store/athletes
	cp -R ./import/ build/
	cp -R ./migrations/ build/

# Run the Dockerfile (only need when working on the dockerfile itself)
run.docker:
	docker run --rm -it -p 3000:3000 robrohan/hungrylegs

# Build all the tools for local playing aroundc
build.all.cli: clean build.cli build.plan

# Create a zip file for the command line tools
distrib.cli: build.all.cli
	mv build hungrylegs
	zip hungrylegs-$(os)-$(hash).zip -r ./hungrylegs/*
	rm -rf hungrylegs

# Create the suffer score file for the dashboard
create_score:
	cd ./cmd/suffer; \
	./suffer.sh "../../store/athletes/UHJvZmVzc29yIFpvb20=.db" ../../score.json
