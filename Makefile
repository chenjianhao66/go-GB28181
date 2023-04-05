BIN_FILE=go-GB28181
RELEASE_FOLDER=release
WIN_APPENDIX=.exe

.PHONY: all
all: check build

.PHONY: setup
setup:
	mkdir -p ${RELEASE_FOLDER}/linux
	mkdir -p ${RELEASE_FOLDER}/darwin
	mkdir -p ${RELEASE_FOLDER}/windows
	cp -r ./config ${RELEASE_FOLDER}/linux
	cp -r ./config ${RELEASE_FOLDER}/darwin
	cp -r ./config ${RELEASE_FOLDER}/windows

.PHONY: build
build: setup
	@GOARCH=amd64 GOOS=linux go build -o "./${RELEASE_FOLDER}/linux/${BIN_FILE}"
	@GOARCH=amd64 GOOS=darwin go build -o "./${RELEASE_FOLDER}/darwin/${BIN_FILE}"
	@GOARCH=amd64 GOOS=windows go build -o "./${RELEASE_FOLDER}/windows/${BIN_FILE}${WIN_APPENDIX}"

.PHONY: clean
clean:
	@rm -rf ./${RELEASE_FOLDER}

.PHONY: test
test:
	@go test

.PHONY: check
check:
	@go fmt ./
	@go vet ./

.PHONY: help
help:
	@echo "make 格式化go代码 并编译生成二进制文件"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"
	@echo "make check 格式化go代码"
	@echo "make run 直接运行程序"
