DIR = $(shell pwd)
BINFILE = $(shell basename ${DIR})
BUILD_DIR = /tmp/gobuild

export GOPATH=$(BUILD_DIR)

build: 
	@mkdir -p $(BUILD_DIR)/src && cp -rf $(DIR) $(BUILD_DIR)/src
	@go build -o $(BINFILE) -ldflags "-w" 
	@rm -rf $(BUILD_DIR)
	@echo "$(BINFILE) is created."

clean:
	@rm $(BINFILE)
