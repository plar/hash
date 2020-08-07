SERVICE_NAME := hashsvc
CUR_DIR      := $(shell pwd)
OUTPUT_DIR   := $(CUR_DIR)/_bin
TEST_DIR     := $(CUR_DIR)/_test
VERSION      := $(shell cat VERSION)

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Build targets:"
	@echo "\tservice            Build the hash service binary"
	@echo ""
	@echo "Test targets:"
	@echo "\ttest               Runs unit tests"
	@echo "\ttest-cover-report  Show the test coverage report, (run make test first)"
	@echo "\tbench              Runs benchmarks"
	@echo ""
	@echo "Docker targets:"
	@echo "\tdocker-service     Build the hash service docker container"
	@echo "\tdocker-run         Run the hash service using docker container, (run make docker-service first)"
	@echo ""
	@echo "Misc targets:"
	@echo "\tclean              Removes build and test artifacts"
	@echo ""

.PHONY: service	
service:
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(SERVICE_NAME) ./cmd/hashsvc/main.go

.PHONY: test
test:
	@mkdir -p $(TEST_DIR)
	@go test -race -coverprofile=$(TEST_DIR)/coverage.out ./...

.PHONY: test-cover-report
test-cover-report:
	@go tool cover -html=$(TEST_DIR)/coverage.out

.PHONY: bench
bench:
	@mkdir -p $(TEST_DIR)
	@go test -benchmem -bench=. -run=^a ./...

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR) $(TEST_DIR)


.PHONE: docker-service
docker-service:
	docker build --build-arg VERSION=$(VERSION) -f Dockerfile -t plar/$(SERVICE_NAME):latest .

.PHONE: docker-run
docker-run:
	docker run -p 8080:8080 plar/hashsvc:latest
