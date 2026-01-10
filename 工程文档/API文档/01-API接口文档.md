# GoAdmin API 接口文档

## 文档信息

- **项目名称**: GoAdmin
- **文档版本**: v1.0
- **创建日期**: 2026-01-06
- **最后更新**: 2026-01-06
- **文档类型**: API 接口文档

---

## 1. API 概述

### 1.1 API 基础信息

**Base URL**: `{prefix}`

**认证方式**: Session 认证

**请求格式**: JSON / Form Data

**响应格式**: JSON / HTML

**字符编码**: UTF-8

### 1.2 URL 前缀配置

GoAdmin 支持自定义 URL 前缀，通过配置文件设置：

```json
{
  "url_prefix": "/admin"
}
```

所有 API 路径都会加上此前缀，例如：
- `/login` → `/admin/login`
- `/info/user` → `/admin/info/user`

### 1.3 URL 格式配置

CRUD 操作的 URL 格式可以通过配置自定义：

```json
{
  "url_format": {
    "info": "/info/:prefix",
    "detail": "/info/:prefix/detail/:id",
    "show_edit": "/info/:prefix/edit/:id",
    "show_new": "/info/:prefix/new",
    "edit": "/info/:prefix/edit",
    "new": "/info/:prefix/new",
    "delete": "/info/:prefix/delete",
    "export": "/info/:prefix/export",
    "update": "/info/:prefix/update"
  }
}
```

### 1.4 认证说明

大部分 API 需要用户登录认证。认证通过 Session 实现：

**登录流程**:
1. 访问登录页面获取 Session
2. 提交用户名和密码
3. 服务器验证并创建 Session
4. 后续请求携带 Session Cookie

**认证中间件**:
- `auth.Middleware(admin.Conn)`: 验证用户登录状态
- `admin.guardian.CheckPrefix`: 验证用户权限

---

## 2. 认证 API

### 2.1 显示登录页面

**接口地址**: `GET /login`

**权限**: 公开

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/login
```

---

### 2.2 用户登录

**接口地址**: `POST /signin`

**权限**: 公开

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| captcha | string | 是 | 验证码 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/signin \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=123456&captcha=1234"
```

**响应**: 重定向到首页或返回错误信息

---

### 2.3 用户登出

**接口地址**: `GET /logout`

**权限**: 需要登录

**请求参数**: 无

**响应**: 重定向到登录页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/logout \
  -H "Cookie: session=xxx"
```

---

### 2.4 服务器登录

**接口地址**: `POST /server/login`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| server | string | 是 | 服务器地址 |
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/server/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "server=192.168.1.100&username=admin&password=123456"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 3. 安装 API

### 3.1 显示安装页面

**接口地址**: `GET /install`

**权限**: 公开

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/install
```

---

### 3.2 检查数据库连接

**接口地址**: `POST /install/database/check`

**权限**: 公开

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| driver | string | 是 | 数据库驱动 (mysql/postgresql/sqlite/mssql) |
| host | string | 是 | 数据库主机 |
| port | string | 是 | 数据库端口 |
| user | string | 是 | 数据库用户名 |
| pwd | string | 是 | 数据库密码 |
| name | string | 是 | 数据库名称 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/install/database/check \
  -H "Content-Type: application/json" \
  -d '{
    "driver": "mysql",
    "host": "localhost",
    "port": "3306",
    "user": "root",
    "pwd": "password",
    "name": "goadmin"
  }'
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "connection": "ok"
  }
}
```

---

## 4. 菜单管理 API

### 4.1 显示菜单列表

**接口地址**: `GET /menu`

**权限**: 需要登录

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/menu \
  -H "Cookie: session=xxx"
```

---

### 4.2 显示新建菜单页面

**接口地址**: `GET /menu/new`

**权限**: 需要登录

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/menu/new \
  -H "Cookie: session=xxx"
```

---

### 4.3 显示编辑菜单页面

**接口地址**: `GET /menu/edit/show`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 菜单 ID |

**响应**: HTML 页面

**示例**:

```bash
curl -X GET "http://localhost:8080/admin/menu/edit/show?id=1" \
  -H "Cookie: session=xxx"
```

---

### 4.4 新建菜单

**接口地址**: `POST /menu/new`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| title | string | 是 | 菜单标题 |
| icon | string | 否 | 菜单图标 |
| uri | string | 否 | 菜单链接 |
| header | string | 否 | 菜单头部 |
| parent_id | int | 否 | 父菜单 ID |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/menu/new \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "title=用户管理&icon=fa-users&uri=/admin/users" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10
  }
}
```

---

### 4.5 编辑菜单

**接口地址**: `POST /menu/edit`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 菜单 ID |
| title | string | 是 | 菜单标题 |
| icon | string | 否 | 菜单图标 |
| uri | string | 否 | 菜单链接 |
| header | string | 否 | 菜单头部 |
| parent_id | int | 否 | 父菜单 ID |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/menu/edit \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=10&title=用户管理&icon=fa-users&uri=/admin/users" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 4.6 删除菜单

**接口地址**: `POST /menu/delete`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 菜单 ID |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/menu/delete \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=10" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 4.7 菜单排序

**接口地址**: `POST /menu/order`

**权限**: 需要登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| ids | string | 是 | 菜单 ID 列表，逗号分隔 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/menu/order \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "ids=1,2,3,4,5" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 5. CRUD 操作 API

### 5.1 显示信息页面

**接口地址**: `GET /info/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/info/user \
  -H "Cookie: session=xxx"
```

**响应**: HTML 页面

---

### 5.2 显示详情页面

**接口地址**: `GET /info/:prefix/detail/:id`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/info/user/detail/1 \
  -H "Cookie: session=xxx"
```

**响应**: HTML 页面

---

### 5.3 显示新建页面

**接口地址**: `GET /info/:prefix/new`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/info/user/new \
  -H "Cookie: session=xxx"
```

**响应**: HTML 页面

---

### 5.4 显示编辑页面

**接口地址**: `GET /info/:prefix/edit/:id`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/info/user/edit/1 \
  -H "Cookie: session=xxx"
```

**响应**: HTML 页面

---

### 5.5 新建记录

**接口地址**: `POST /info/:prefix/new`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**: 根据表结构动态生成

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/info/user/new \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=testuser&email=test@example.com&password=123456" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10
  }
}
```

---

### 5.6 编辑记录

**接口地址**: `POST /info/:prefix/edit`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |
| 其他字段 | - | - | 根据表结构动态生成 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/info/user/edit \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1&username=admin&email=admin@example.com" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 5.7 删除记录

**接口地址**: `POST /info/:prefix/delete`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/info/user/delete \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 5.8 更新记录

**接口地址**: `POST /info/:prefix/update`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |
| 其他字段 | - | - | 根据表结构动态生成 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/info/user/update \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1&username=admin&email=admin@example.com" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 5.9 导出记录

**接口地址**: `POST /info/:prefix/export`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| export_type | string | 否 | 导出类型 (csv/xlsx) |
| ids | string | 否 | 要导出的 ID 列表，逗号分隔 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/info/user/export \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "export_type=xlsx" \
  -H "Cookie: session=xxx" \
  -o users.xlsx
```

**响应**: 文件下载

---

## 6. JSON API

JSON API 需要在配置中启用：

```json
{
  "open_admin_api": true
}
```

### 6.1 获取列表

**接口地址**: `GET /api/list/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 10 |
| sortField | string | 否 | 排序字段 |
| sortType | string | 否 | 排序类型 (asc/desc) |
| 其他字段 | - | - | 过滤条件 |

**请求示例**:

```bash
curl -X GET "http://localhost:8080/admin/api/list/user?page=1&pageSize=10" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "list": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com"
      }
    ]
  }
}
```

---

### 6.2 获取详情

**接口地址**: `GET /api/detail/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X GET "http://localhost:8080/admin/api/detail/user?id=1" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }
}
```

---

### 6.3 获取创建表单

**接口地址**: `GET /api/create/form/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/api/create/form/user \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "fields": [
      {
        "name": "username",
        "type": "text",
        "label": "用户名",
        "required": true
      },
      {
        "name": "email",
        "type": "email",
        "label": "邮箱",
        "required": true
      }
    ]
  }
}
```

---

### 6.4 创建记录

**接口地址**: `POST /api/create/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**: 根据表结构动态生成

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/api/create/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }' \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10
  }
}
```

---

### 6.5 获取编辑表单

**接口地址**: `GET /api/edit/form/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X GET "http://localhost:8080/admin/api/edit/form/user?id=1" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "fields": [
      {
        "name": "username",
        "type": "text",
        "label": "用户名",
        "value": "admin",
        "required": true
      },
      {
        "name": "email",
        "type": "email",
        "label": "邮箱",
        "value": "admin@example.com",
        "required": true
      }
    ]
  }
}
```

---

### 6.6 编辑记录

**接口地址**: `POST /api/edit/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |
| 其他字段 | - | - | 根据表结构动态生成 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/api/edit/user \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }' \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 6.7 删除记录

**接口地址**: `POST /api/delete/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/api/delete/user \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 6.8 更新记录

**接口地址**: `POST /api/update/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| id | int | 是 | 记录 ID |
| 其他字段 | - | - | 根据表结构动态生成 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/api/update/user \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }' \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 6.9 导出记录

**接口地址**: `POST /api/export/:prefix`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| prefix | string | 是 | 表前缀 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| export_type | string | 否 | 导出类型 (csv/xlsx) |
| ids | string | 否 | 要导出的 ID 列表，逗号分隔 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/api/export/user \
  -H "Content-Type: application/json" \
  -d '{
    "export_type": "xlsx"
  }' \
  -H "Cookie: session=xxx" \
  -o users.xlsx
```

**响应**: 文件下载

---

## 7. 插件管理 API

### 7.1 显示插件列表

**接口地址**: `GET /plugins`

**权限**: 需要登录

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/plugins \
  -H "Cookie: session=xxx"
```

---

### 7.2 显示插件商店（仅非生产环境）

**接口地址**: `GET /plugins/store`

**权限**: 需要登录

**环境要求**: 非生产环境

**请求参数**: 无

**响应**: HTML 页面

**示例**:

```bash
curl -X GET http://localhost:8080/admin/plugins/store \
  -H "Cookie: session=xxx"
```

---

### 7.3 下载插件（仅非生产环境）

**接口地址**: `POST /plugin/download`

**权限**: 需要登录

**环境要求**: 非生产环境

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| name | string | 是 | 插件名称 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/plugin/download \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=example-plugin" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 7.4 获取插件详情（仅非生产环境）

**接口地址**: `POST /plugin/detail`

**权限**: 需要登录

**环境要求**: 非生产环境

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| name | string | 是 | 插件名称 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/plugin/detail \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=example-plugin" \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "name": "example-plugin",
    "version": "1.0.0",
    "description": "Example plugin",
    "author": "author"
  }
}
```

---

## 8. 系统信息 API

### 8.1 获取系统信息

**接口地址**: `GET /application/info`

**权限**: 需要登录

**请求参数**: 无

**响应**: JSON

**示例**:

```bash
curl -X GET http://localhost:8080/admin/application/info \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "go_version": "go1.24.2",
    "os": "linux",
    "arch": "amd64",
    "uptime": "24h30m15s"
  }
}
```

---

## 9. 操作 API

### 9.1 执行操作

**接口地址**: `ANY /operation/:__goadmin_op_id`

**权限**: 需要登录

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| __goadmin_op_id | string | 是 | 操作 ID |

**请求参数**: 根据操作类型动态生成

**请求示例**:

```bash
curl -X POST http://localhost:8080/admin/operation/custom_operation \
  -H "Content-Type: application/json" \
  -d '{
    "param1": "value1",
    "param2": "value2"
  }' \
  -H "Cookie: session=xxx"
```

**响应**: JSON

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 10. 静态资源 API

### 10.1 获取静态资源

**接口地址**: `GET /assets/*`

**权限**: 公开

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| * | string | 是 | 资源路径 |

**请求示例**:

```bash
curl -X GET http://localhost:8080/admin/assets/css/style.css
```

**响应**: 静态文件

---

## 11. 错误码说明

### 11.1 通用错误码

| 错误码 | 说明 |
|-------|------|
| 0 | 成功 |
| 1000 | 未知错误 |
| 1001 | 参数错误 |
| 1002 | 数据库错误 |
| 1003 | 权限不足 |
| 1004 | 未登录 |
| 1005 | 数据不存在 |
| 1006 | 数据已存在 |

### 11.2 认证错误码

| 错误码 | 说明 |
|-------|------|
| 2000 | 用户名或密码错误 |
| 2001 | 验证码错误 |
| 2002 | 账户已被禁用 |
| 2003 | Session 已过期 |

### 11.3 数据库错误码

| 错误码 | 说明 |
|-------|------|
| 3000 | 数据库连接失败 |
| 3001 | 表不存在 |
| 3002 | 字段不存在 |
| 3003 | 数据类型错误 |
| 3004 | 约束冲突 |

### 11.4 文件错误码

| 错误码 | 说明 |
|-------|------|
| 4000 | 文件不存在 |
| 4001 | 文件类型不支持 |
| 4002 | 文件大小超限 |
| 4003 | 文件上传失败 |

---

## 12. 响应格式

### 12.1 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 12.2 错误响应

```json
{
  "code": 1001,
  "message": "参数错误",
  "data": {}
}
```

### 12.3 列表响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "list": []
  }
}
```

---

## 13. 请求头说明

### 13.1 通用请求头

| 请求头 | 说明 | 示例 |
|-------|------|------|
| Content-Type | 请求内容类型 | application/json |
| Cookie | Session Cookie | session=xxx |
| User-Agent | 用户代理 | Mozilla/5.0 |

### 13.2 响应头

| 响应头 | 说明 | 示例 |
|-------|------|------|
| Content-Type | 响应内容类型 | application/json |
| x-request-id | 请求 ID | 1234567890 |

---

## 14. 使用示例

### 14.1 完整的 CRUD 流程

```bash
# 1. 登录
curl -X POST http://localhost:8080/admin/signin \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=123456&captcha=1234" \
  -c cookies.txt

# 2. 获取用户列表
curl -X GET "http://localhost:8080/admin/api/list/user?page=1&pageSize=10" \
  -b cookies.txt

# 3. 获取创建表单
curl -X GET http://localhost:8080/admin/api/create/form/user \
  -b cookies.txt

# 4. 创建用户
curl -X POST http://localhost:8080/admin/api/create/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }' \
  -b cookies.txt

# 5. 获取编辑表单
curl -X GET "http://localhost:8080/admin/api/edit/form/user?id=1" \
  -b cookies.txt

# 6. 编辑用户
curl -X POST http://localhost:8080/admin/api/edit/user \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }' \
  -b cookies.txt

# 7. 删除用户
curl -X POST http://localhost:8080/admin/api/delete/user \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  -b cookies.txt

# 8. 登出
curl -X GET http://localhost:8080/admin/logout \
  -b cookies.txt
```

### 14.2 使用 JavaScript 调用 API

```javascript
// 登录
async function login(username, password, captcha) {
  const response = await fetch('/admin/signin', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: `username=${username}&password=${password}&captcha=${captcha}`,
  });
  return response.json();
}

// 获取列表
async function getList(prefix, page = 1, pageSize = 10) {
  const response = await fetch(`/admin/api/list/${prefix}?page=${page}&pageSize=${pageSize}`, {
    method: 'GET',
    credentials: 'include',
  });
  return response.json();
}

// 创建记录
async function create(prefix, data) {
  const response = await fetch(`/admin/api/create/${prefix}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(data),
  });
  return response.json();
}

// 编辑记录
async function update(prefix, id, data) {
  const response = await fetch(`/admin/api/edit/${prefix}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ id, ...data }),
  });
  return response.json();
}

// 删除记录
async function delete(prefix, id) {
  const response = await fetch(`/admin/api/delete/${prefix}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ id }),
  });
  return response.json();
}

// 使用示例
(async () => {
  // 登录
  await login('admin', '123456', '1234');

  // 获取用户列表
  const list = await getList('user', 1, 10);
  console.log(list);

  // 创建用户
  const created = await create('user', {
    username: 'testuser',
    email: 'test@example.com',
    password: '123456',
  });
  console.log(created);

  // 编辑用户
  const updated = await update('user', 1, {
    username: 'admin',
    email: 'admin@example.com',
  });
  console.log(updated);

  // 删除用户
  const deleted = await delete('user', 1);
  console.log(deleted);
})();
```

---

## 15. 注意事项

### 15.1 安全性

1. **HTTPS**: 生产环境必须使用 HTTPS
2. **CSRF 防护**: 所有 POST 请求需要 CSRF Token
3. **输入验证**: 所有输入参数必须验证
4. **SQL 注入**: 使用参数化查询，避免 SQL 注入
5. **XSS 防护**: 输出时进行 HTML 转义

### 15.2 性能优化

1. **分页查询**: 使用分页避免大量数据查询
2. **索引优化**: 为常用查询字段添加索引
3. **缓存策略**: 使用缓存减少数据库查询
4. **连接池**: 合理配置数据库连接池

### 15.3 错误处理

1. **统一错误码**: 使用统一的错误码体系
2. **错误日志**: 记录所有错误信息
3. **用户友好**: 向用户返回友好的错误提示
4. **错误追踪**: 使用 Trace ID 追踪错误

---

## 16. 版本历史

| 版本 | 日期 | 变更说明 |
|-----|------|---------|
| v1.0 | 2026-01-06 | 初始版本，包含完整的 API 文档 |

---

## 17. 附录

### 17.1 术语表

| 术语 | 说明 |
|-----|------|
| prefix | 表前缀，用于标识不同的数据表 |
| Session | 用户会话，用于认证和状态管理 |
| CSRF Token | 跨站请求伪造令牌，用于防止 CSRF 攻击 |
| Trace ID | 请求追踪 ID，用于追踪请求链路 |

### 17.2 参考资源

- **GoAdmin 官方文档**: http://doc.go-admin.cn/zh
- **GoAdmin GitHub**: https://github.com/purpose168/GoAdmin
- **GoAdmin Demo**: https://demo.go-admin.cn

### 17.3 联系方式

- **QQ 群**: 694446792
- **Discord**: https://discord.gg/usAaEpCP
- **Telegram**: https://t.me/joinchat/NlyH6Bch2QARZkArithKvg
- **GitHub Issues**: https://github.com/purpose168/GoAdmin/issues

---

## 总结

GoAdmin API 提供了完整的 CRUD 操作、认证授权、菜单管理、插件管理等功能。通过本文档，开发者可以快速了解和使用 GoAdmin 的 API 接口。

API 的主要特点：

1. **RESTful 风格**: 遵循 RESTful API 设计规范
2. **统一响应格式**: 所有 API 返回统一的 JSON 格式
3. **完善的认证**: 基于 Session 的认证机制
4. **灵活的权限**: 基于 RBAC 的权限控制
5. **丰富的功能**: 支持 CRUD、导出、操作等多种功能

如有任何问题或建议，欢迎通过上述联系方式反馈。
