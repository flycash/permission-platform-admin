# 初始化项目环境
.PHONY: setup
setup:
	@sh ./scripts/setup.sh

# 格式化代码
.PHONY: fmt
fmt:
	@goimports -l -w $$(find . -type f -name '*.go' -not -path "./.idea/*" -not -path "./**/ioc/wire_gen.go" -not -path "./**/ioc/wire.go")
	@gofumpt -l -w $$(find . -type f -name '*.go' -not -path "./.idea/*" -not -path "./**/ioc/wire_gen.go" -not -path "./**/ioc/wire.go")

# 清理项目依赖
.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) --no-print-directory fmt
	@$(MAKE) --no-print-directory tidy

# 代码规范检查
.PHONY: lint
lint:
	@golangci-lint run -c ./scripts/lint/.golangci.yaml ./...

.PHONY: run
run:
	@export EGO_DEBUG=true && go run main.go --config=./config/config.yaml