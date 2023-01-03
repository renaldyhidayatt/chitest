# run: as running app
run:
	go run main.go

## test: Run all test in this app
test:
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test: Run all test with clean cache in this app
test_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test_cover: Run all test with coverage
test_cover:
	@echo "All test are running with coverage..."
	go test ./... -v -cover

## test: Run all test with clean cache and coverage
test_cover_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test ./... -v -cover
	@echo "Test finished"
