# Go settings
GO        := go
GOOS      ?= linux
GOARCH    ?= amd64
BINARY    ?= ufWall
BIN_DIR   := bin
SRC       := ./cmd/ufWall
LDFLAGS   := -s -w

.PHONY: all help build run clean dist

all: build

help:
	@echo "ufWall Makefile Utilities"
	@echo "  help      -- Show this message (default)"
	@echo "  build     -- Build binary into $(BIN_DIR)/$(BINARY)"
	@echo "  run       -- Run the binary (sudo required, uses $(BIN_DIR)/$(BINARY))"
	@echo "  clean     -- Remove the ./bin directory"
	@echo "  dist      -- Create a .tar.gz with a statically‑linked binary"

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build \
		-ldflags="$(LDFLAGS)" \
		-o $(BIN_DIR)/$(BINARY) \
		$(SRC)

run: build
	@sudo $(BIN_DIR)/$(BINARY)

clean:
	@rm -rf $(BIN_DIR) dist

dist: clean build
	@mkdir -p dist
	@cp $(BIN_DIR)/$(BINARY) dist/
	@cp README.md dist/
	@tar czf dist/$(BINARY)-$(GOOS)-$(GOARCH).tar.gz -C dist $(BINARY) README.md
	@echo "Distribution packaged in dist/$(BINARY)-$(GOOS)-$(GOARCH).tar.gz"
