

LAMBDAS = $(sort $(dir $(wildcard cmd/renwick/*/)))
BUILD_DIR=artifacts


all: ensure $(LAMBDAS)
	@echo $@  $(LAMBDAS)

ensure:
	dep ensure

clean:
	@[ ! -z "$(BUILD_DIR)" ] && rm -rf $(BUILD_DIR)

 $(LAMBDAS):
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$@lambda $@main.go


.PHONY: all ensure $(LAMBDAS)