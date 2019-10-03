NAME    := dstcluster
GOARCH  := amd64
TARGETS := windows linux darwin
OUTPUT  := .build
PACKAGE  = $(NAME)-$(GOOS)-$(GOARCH)
BINARY   = $(NAME)$(EXTENSION)

PACKAGE_PATH = $(OUTPUT)/$(PACKAGE)
BINARY_PATH  = $(PACKAGE_PATH)/$(BINARY)


.PHONY: all clean $(TARGETS)

# Build all targets
all: $(TARGETS)

# Clean build directory
clean:
	rm -rf .build/

# Export GOOS for every target
$(TARGETS): export GOOS = $@

# Build for windows platform
windows: EXTENSION = .exe
windows:
	$(MAKE) $(PACKAGE_PATH).zip

# Build for unix platform
linux darwin:
	$(MAKE) $(PACKAGE_PATH).tar.gz

# Create zip package from binary
$(PACKAGE_PATH).zip: $(BINARY_PATH)
	cd $(PACKAGE_PATH) && zip ../$(PACKAGE).zip $(BINARY)

# Create gzipped tar package from binary
$(PACKAGE_PATH).tar.gz: $(BINARY_PATH)
	cd $(PACKAGE_PATH) && tar czf ../$(PACKAGE).tar.gz $(BINARY)

# Build binary
$(BINARY_PATH):
	mkdir -p $(dir PACKAGE_PATH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_PATH) .
