COVERAGE_FILE=/tmp/abutils_coverage.out

test:
	go test -cover -v -race

cover:
	go test -coverprofile=$(COVERAGE_FILE) && \
	go tool cover -html=$(COVERAGE_FILE) && \
	rm $(COVERAGE_FILE)
