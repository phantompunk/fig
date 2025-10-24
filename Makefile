.PHONY: help dev install clean

BINARY = fig
AIR_BINARY = ./tmp/fig
INSTALL_PATH=$(HOME)/go/bin

# Set the default goal
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Usage:\n  make [command]"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## dev: Live reloading for cli using air
dev: install
	@air

## install: Symlink dev cli to $GOPATH/bin
install:
	@mkdir -p $(INSTALL_PATH)
	@ln -sf $(PWD)/$(AIR_BINARY) $(INSTALL_PATH)/$(BINARY)
	@echo "Symlinked $(INSTALL_PATH)/$(BINARY) -> $(PWD)/$(AIR_BINARY)"

## clean: Remove all development artifacts
clean:
	@rm -f $(INSTALL_PATH)/$(BINARY)
	@echo "Removed symlink $(INSTALL_PATH)/$(BINARY)"
