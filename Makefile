# Go 命令行工具
GOCMD = go
# Go 构建命令
GOBUILD = $(GOCMD) build

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

# 默认目标：执行测试
all: test

## 测试相关命令 (tests)

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

## 代码风格检查 (code style check)

# 执行所有代码检查：格式化、golint、govet、cilint
lint: fmt golint govet cilint

# 格式化代码
fmt:
	# 使用 go fmt 格式化代码
	GO111MODULE=off go fmt ./...
	# 使用 goimports 格式化并整理导入
	GO111MODULE=off goimports -l -w .

# 执行 go vet 静态分析
govet:
	# 运行 go vet 检查代码问题
	GO111MODULE=off go vet ./...

# 执行 golangci-lint 检查
cilint:
	# 运行 golangci-lint 进行综合代码检查
	GO111MODULE=off golangci-lint run

# 执行 golint 检查
golint:
	# 运行 golint 检查代码风格
	GO111MODULE=off golint ./...

# 构建模板文件
build-tmpl:
    ## 编译表单模板 (form tmpl build)
	# 将模板目录编译为 Go 代码
	adm compile tpl --src ./template/types/tmpls/ --dist ./template/types/tmpl.go --package types --var tmpls
    ## 编译生成器模板 (generator tmpl build)
	# 将表格模板目录编译为 Go 代码
	adm compile tpl --src ./plugins/admin/modules/table/tmpl --dist ./plugins/admin/modules/table/tmpl.go --package table --var tmpls

# 声明伪目标，避免与同名文件冲突
.PHONY: all fmt golint govet cp-mod restore-mod test black-box-test mysql-test sqlite-test import-sqlite import-mysql import-postgresql pg-test lint cilint cli
