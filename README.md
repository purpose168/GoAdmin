<p align="center">
  <a href="https://github.com/GoAdminGroup/go-admin">
    <img width="48%" alt="go-admin" src="http://quick.go-admin.cn/official/assets/imgs/github_logo.png">
  </a>
</p>

<p align="center">
    缺失的 Golang 数据管理面板构建工具。
</p>

<p align="center">
    <a href="https://book.go-admin.cn/en">英文文档</a> |
	<a href="http://doc.go-admin.cn/zh/">中文文档</a> |
    <a href="./README_CN.md">中文介绍</a> |
    <a href="https://demo.go-admin.com">英文演示</a> |
    <a href="https://demo.go-admin.cn">中文演示</a> |
    <a href="https://twitter.com/cg3365688034">Twitter</a> |
    <a href="http://discuss.go-admin.com">论坛</a>
</p>

<p align="center">
  <a href="http://drone.go-admin.com/GoAdminGroup/go-admin"><img alt="构建状态" src="http://drone.go-admin.com/api/badges/purpose168/GoAdmin/status.svg?ref=refs/heads/master"></a>
  <a href="https://goreportcard.com/report/github.com/GoAdminGroup/go-admin"><img alt="Go 报告卡" src="https://goreportcard.com/badge/github.com/GoAdminGroup/go-admin"></a>
  <a href="https://goreportcard.com/report/github.com/GoAdminGroup/go-admin"><img alt="golang" src="https://img.shields.io/badge/awesome-golang-blue.svg"></a>
  <a href="https://discord.gg/usAaEpCP"><img alt="discord" src="https://img.shields.io/badge/chat%20on-Discord-blue.svg"></a>
  <a href="https://t.me/joinchat/NlyH6Bch2QARZkArithKvg" rel="nofollow"><img alt="telegram" src="https://img.shields.io/badge/chat%20on-telegram-blue" style="max-width:100%;"></a>
  <a href="https://raw.githubusercontent.com/purpose168/GoAdmin/master/LICENSE" rel="nofollow"><img src="https://img.shields.io/badge/license-Apache2.0-blue.svg" alt="许可证" data-canonical-src="https://img.shields.io/badge/license-Apache2.0-blue.svg" style="max-width:100%;"></a>
</p>

<p align="center">
    灵感来源于 <a href="https://github.com/z-song/laravel-admin" target="_blank">laravel-admin</a>
</p>

## 前言

GoAdmin 是一个工具包，帮助您为 Golang 应用程序构建数据可视化管理面板。

在线演示：[https://demo.go-admin.com](https://demo.go-admin.com)

![界面](http://file.go-admin.cn/introduction/interface_en_3.png)

## 特性

- 🚀 **快速**：在 **十** 分钟内构建一个生产级管理面板应用。
- 🎨 **主题**：支持美观的 UI 主题（默认使用 adminlte，更多主题即将推出）。
- 🔢 **插件**：提供许多可用的插件（更多实用且强大的插件即将推出）。
- ✅ **权限控制**：开箱即用的 RBAC（基于角色的访问控制）认证系统。
- ⚙️ **框架**：支持大多数 Go Web 框架。

## 翻译

我们需要您的帮助：[https://github.com/GoAdminGroup/docs/issues/1](https://github.com/GoAdminGroup/docs/issues/1)

## 谁在使用

[在 issue 中评论告诉我们](https://github.com/purpose168/GoAdmin/issues/71)。

## 如何使用

按照以下三个步骤运行。

```shell
# 创建新项目目录并进入该目录
$ mkdir new_project && cd new_project
# 安装 adm 命令行工具
$ go install github.com/GoAdminGroup/adm@latest
# 初始化 web 项目
$ adm init web
```

## 示例

快速跟进示例：

- [纯 Golang](https://github.com/GoAdminGroup/example)，简单且依赖较少
- [带前端模板的 Golang](https://github.com/GoAdminGroup/example_with_frontend)，可自行更改模板
- [带 Vue 的 Golang](https://github.com/GoAdminGroup/example_with_vue)，如果您有 Vue 经验

查看 [文档](https://book.go-admin.cn) 了解更多详情。

## 赞助者

您的支持将帮助我做得更好！[[成为赞助者](https://opencollective.com/go-admin#backer)]
 <a href="https://opencollective.com/go-admin#backers" target="_blank"><img src="https://opencollective.com/go-admin/backers.svg?width=890"></a>

## 贡献

[点击查看贡献指南](CONTRIBUTING.md)

<strong>点击此处加入开发团队</strong>

[加入 Telegram](https://t.me/joinchat/NlyH6Bch2QARZkArithKvg)
