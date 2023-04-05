.PHONY: all setup build clean run check help
BIN_FILE=go-GB28181

all: check build

setup:
	mkdir -p build/linux
	mkdir -p build/darwin
	mkdir -p build/windows

build: setup
	@GOARCH=amd64 GOOS=linux go build -o "./build/linux/${BIN_FILE}"
	@GOARCH=amd64 GOOS=darwin go build -o "./build/darwin/${BIN_FILE}"
	@GOARCH=amd64 GOOS=windows go build -o "./build/windows/${BIN_FILE}"

clean:
	@rm -rf ./build

test:
	@go test

check:
	@go fmt ./
	@go vet ./

run:
	./"${BIN_FILE}"

help:
	@echo "make 格式化go代码 并编译生成二进制文件"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"
	@echo "make check 格式化go代码"
	@echo "make run 直接运行程序"