

LAMBDAS = $(sort $(dir $(wildcard cmd/*/)))
BUILD_DIR=artifacts


all: ensure $(LAMBDAS)
	@echo $@  $(LAMBDAS)

ensure:
	dep ensure

clean:
	@[ ! -z "$(BUILD_DIR)" ] && rm -rf $(BUILD_DIR)

 $(LAMBDAS):
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$@lambda $@main.go
	@(cd $(BUILD_DIR)/$@; zip lambda.zip ./lambda;)


.PHONY: all ensure $(LAMBDAS)