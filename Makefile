# Find all .go files in the directory and its subdirectories
GO_FILES := $(shell find . -name '*.go')

# Target for building the gonion binary
gonion: $(GO_FILES)
	@go build -o ./gonion ./cmd/gonion

# Target for running all tests
test:
	@go test -v ./...
