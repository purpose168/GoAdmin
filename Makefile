# 定义 Go 命令工具
GOCMD = go
# 定义构建命令
GOBUILD = $(GOCMD) build
# 定义安装命令
GOINSTALL = $(GOCMD) install
# 定义测试命令
GOTEST = $(GOCMD) test
# 定义生成的二进制文件名称
BINARY_NAME = goadmin
# 定义 CLI 工具名称
CLI = adm

# 测试配置文件路径 - MySQL 配置
TEST_CONFIG_PATH=./../../common/config.json
# 测试配置文件路径 - PostgreSQL 配置
TEST_CONFIG_PQ_PATH=./../../common/config_pg.json
# 测试配置文件路径 - SQLite 配置
TEST_CONFIG_SQLITE_PATH=./../../common/config_sqlite.json
# 测试配置文件路径 - SQL Server 配置
TEST_CONFIG_MS_PATH=./../../common/config_ms.json
# 测试框架目录
TEST_FRAMEWORK_DIR=./tests/frameworks

## 数据库配置 (database configs)
# MySQL 数据库主机地址
MYSQL_HOST = db_mysql
# MySQL 数据库端口
MYSQL_PORT = 3306
# MySQL 数据库用户名
MYSQL_USER = root
# MySQL 数据库密码
MYSQL_PWD = root

# PostgreSQL 数据库主机地址
POSTGRESSQL_HOST = db_pgsql
# PostgreSQL 数据库端口
POSTGRESSQL_PORT = 5432
# PostgreSQL 数据库用户名
POSTGRESSQL_USER = postgres
# PostgreSQL 数据库密码
POSTGRESSQL_PWD = root

# 测试数据库名称
TEST_DB = go-admin-test

# 默认目标：启动服务
all: serve

# ------------------------
# 服务管理
# ------------------------

# 启动服务：直接运行当前目录的 Go 程序
serve:
	@echo "=== 启动服务 ==="
	$(GOCMD) run .

# ------------------------
# 构建管理
# ------------------------

# 构建项目：生成 Linux 平台的二进制文件
build:
	@echo "=== 构建项目 ==="
	@mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_NAME) -v ./

# ------------------------
# 依赖管理
# ------------------------

# 清理模块缓存：清理 Go 模块缓存，解决依赖冲突
mod-clean:
	@echo "=== 清理模块缓存 ==="
	$(GOCMD) clean -modcache

# 修复依赖关系：更新 go.mod 和 go.sum 文件，确保依赖关系正确
mod-tidy:
	@echo "=== 修复依赖关系 ==="
	$(GOCMD) mod tidy

# 生成 vendor 目录：将依赖复制到本地 vendor 目录
mod-vendor:
	@echo "=== 生成 vendor 目录 ==="
	$(GOCMD) mod vendor

# 检查依赖关系：检查依赖关系是否正确
mod-verify:
	@echo "=== 检查依赖关系 ==="
	$(GOCMD) mod verify

# 调试依赖问题：打印依赖图，用于调试依赖问题
mod-graph:
	@echo "=== 打印依赖图 ==="
	$(GOCMD) mod graph

# 更新依赖：更新所有依赖到最新版本
mod-update:
	@echo "=== 更新依赖 ==="
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

# ------------------------
# 测试管理
# ------------------------

# 执行完整测试：备份模块 -> 黑盒测试 -> Web 测试 -> 恢复模块
test: cp-mod black-box-test web-test restore-mod

## 测试：黑盒测试 (tests: black box tests)

# 执行所有数据库的黑盒测试
black-box-test: mysql-test pg-test sqlite-test ms-test

# MySQL 数据库黑盒测试
mysql-test: $(TEST_FRAMEWORK_DIR)/*
	# 移除 ugorji/go/codec 依赖以避免冲突
	go get github.com/ugorji/go/codec@none
	# 遍历所有测试框架目录
	for file in $^ ; do \
	# 导入 MySQL 测试数据
	make import-mysql ; \
	# 执行测试：使用模块模式、禁用内联优化、详细输出
	go test -mod=mod -gcflags=all=-l -v ./$${file}/... -args $(TEST_CONFIG_PATH) ; \
	done

# SQLite 数据库黑盒测试
sqlite-test: $(TEST_FRAMEWORK_DIR)/*
	# 遍历所有测试框架目录
	for file in $^ ; do \
	# 导入 SQLite 测试数据
	make import-sqlite ; \
	# 执行测试：使用模块模式、禁用内联优化
	go test -mod=mod -gcflags=all=-l ./$${file}/... -args $(TEST_CONFIG_SQLITE_PATH) ; \
	done

# PostgreSQL 数据库黑盒测试
pg-test: $(TEST_FRAMEWORK_DIR)/*
	# 遍历所有测试框架目录
	for file in $^ ; do \
	# 导入 PostgreSQL 测试数据
	make import-postgresql ; \
	# 执行测试：使用模块模式、禁用内联优化
	go test -mod=mod -gcflags=all=-l ./$${file}/... -args $(TEST_CONFIG_PQ_PATH) ; \
	done

# SQL Server 数据库黑盒测试
ms-test: $(TEST_FRAMEWORK_DIR)/*
	# 遍历所有测试框架目录
	for file in $^ ; do \
	# 导入 SQL Server 测试数据
	make import-mssql ; \
	# 执行测试：使用模块模式、禁用内联优化
	go test -mod=mod -gcflags=all=-l ./$${file}/... -args $(TEST_CONFIG_MS_PATH) ; \
	done

## 测试：用户验收测试 (tests: user acceptance tests)

# Web 端用户验收测试
web-test: import-mysql
	# 执行 Web 测试
	go test -mod=mod ./tests/web/...
	# 清理测试生成的用户文件
	rm -rf ./tests/web/User*

# Web 端调试模式测试
web-test-debug: import-mysql
	# 执行 Web 测试并启用调试模式
	go test -mod=mod ./tests/web/... -args true

## 测试：单元测试 (tests: unit tests)

# 执行单元测试
unit-test:
	# 测试 adm 模块
	go test -mod=mod ./adm/...
	# 测试 context 模块
	go test -mod=mod ./context/...
	# 测试 modules 模块
	go test -mod=mod ./modules/...
	# 测试 admin controller 模块
	go test -mod=mod ./plugins/admin/controller/...
	# 测试 admin parameter 模块
	go test -mod=mod ./plugins/admin/modules/parameter/...
	# 测试 admin table 模块
	go test -mod=mod ./plugins/admin/modules/table/...
	# 测试 admin modules 模块
	go test -mod=mod ./plugins/admin/modules/...

## 测试：辅助命令 (tests: helpers)

# 导入 SQLite 测试数据
import-sqlite:
	# 删除旧的 SQLite 数据库文件
	rm -rf ./tests/common/admin.db
	# 复制测试数据到指定位置
	cp ./tests/data/admin.db ./tests/common/admin.db

# 导入 MySQL 测试数据
import-mysql:
	# 创建测试数据库（如果不存在）
	mysql -h$(MYSQL_HOST) -P${MYSQL_PORT} -u${MYSQL_USER} -p${MYSQL_PWD} -e "create database if not exists \`${TEST_DB}\`"
	# 导入 SQL 脚本到测试数据库
	mysql -h$(MYSQL_HOST) -P${MYSQL_PORT} -u${MYSQL_USER} -p${MYSQL_PWD} ${TEST_DB} < ./tests/data/admin.sql

# 导入 PostgreSQL 测试数据
import-postgresql:
	# 删除已存在的测试数据库
	PGPASSWORD=${POSTGRESSQL_PWD} dropdb -h ${POSTGRESSQL_HOST} -p ${POSTGRESSQL_PORT} -U ${POSTGRESSQL_USER} ${TEST_DB}
	# 创建新的测试数据库
	PGPASSWORD=${POSTGRESSQL_PWD} createdb -h ${POSTGRESSQL_HOST} -p ${POSTGRESSQL_PORT} -U ${POSTGRESSQL_USER} ${TEST_DB}
	# 导入 SQL 脚本到测试数据库
	PGPASSWORD=${POSTGRESSQL_PWD} psql -h ${POSTGRESSQL_HOST} -p ${POSTGRESSQL_PORT} -d ${TEST_DB} -U ${POSTGRESSQL_USER} -f ./tests/data/admin_pg.sql

# 导入 SQL Server 测试数据
import-mssql:
	# 从备份文件恢复 SQL Server 数据库
	/opt/mssql-tools/bin/sqlcmd -S db_mssql -U SA -P Aa123456 -Q "RESTORE DATABASE [goadmin] FROM DISK = N'/home/data/admin_ms.bak' WITH FILE = 1, NOUNLOAD, REPLACE, RECOVERY, STATS = 5"

# 备份 SQL Server 数据库
backup-mssql:
	# 将 SQL Server 数据库备份到指定文件
	docker exec mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P Aa123456 -Q "BACKUP DATABASE [goadmin] TO DISK = N'/home/data/admin_ms.bak' WITH NOFORMAT, NOINIT, NAME = 'goadmin-full', SKIP, NOREWIND, NOUNLOAD, STATS = 10"

# 备份 Go 模块文件
cp-mod:
	# 备份 go.mod 文件
	cp go.mod go.mod.old
	# 备份 go.sum 文件
	cp go.sum go.sum.old

# 恢复 Go 模块文件
restore-mod:
	# 恢复 go.mod 文件
	mv go.mod.old go.mod
	# 恢复 go.sum 文件
	mv go.sum.old go.sum

# 准备测试数据：复制数据库文件用于测试
ready-for-data:
	@echo "=== 准备测试数据 ==="
	@cp admin.db admin_test.db

# 清理测试数据：删除测试数据库文件
clean:
	@echo "=== 清理测试数据 ==="
	@rm -f admin_test.db

# ------------------------
# 代码生成
# ------------------------

# 生成代码：安装 go-admin CLI 工具并生成代码
generate:
	@echo "=== 生成代码 ==="
	$(GOINSTALL) github.com/purpose168/GoAdmin-adm
	$(CLI) generate -c adm_config.ini

# ------------------------
# 开发辅助
# ------------------------

# 执行所有代码检查：格式化、golint、govet、cilint
lint: fmt golint govet cilint

# 格式化代码
fmt:
	@echo "=== 格式化代码 ==="
	$(GOCMD) fmt ./...

# 执行 go vet 静态分析
govet:
	@echo "=== 检查代码 ==="
	$(GOCMD) vet ./...

# 执行 golangci-lint 检查
cilint:
	@echo "=== 执行 golangci-lint 检查 ==="
	golangci-lint run

# 执行 golint 检查
golint:
	@echo "=== 执行 golint 检查 ==="
	golint ./...

# 静态分析：使用 staticcheck 进行静态分析
staticcheck:
	@echo "=== 静态分析 ==="
	@which staticcheck > /dev/null 2>&1 || $(GOCMD) install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# 构建模板文件
build-tmpl:
	## 编译表单模板 (form tmpl build)
	@echo "=== 编译表单模板 ==="
	# 将模板目录编译为 Go 代码
	$(CLI) compile tpl --src ./template/types/tmpls/ --dist ./template/types/tmpl.go --package types --var tmpls
	## 编译生成器模板 (generator tmpl build)
	@echo "=== 编译生成器模板 ==="
	# 将表格模板目录编译为 Go 代码
	$(CLI) compile tpl --src ./plugins/admin/modules/table/tmpl --dist ./plugins/admin/modules/table/tmpl.go --package table --var tmpls

# ------------------------
# 声明伪目标：这些目标不代表实际文件
# ------------------------

.PHONY: all serve build \
	mod-clean mod-tidy mod-vendor mod-verify mod-graph mod-update \
	test black-box-test web-test web-test-debug unit-test mysql-test pg-test sqlite-test ms-test \
	import-sqlite import-mysql import-postgresql import-mssql backup-mssql cp-mod restore-mod ready-for-data clean \
	generate fmt golint govet cilint staticcheck build-tmpl
