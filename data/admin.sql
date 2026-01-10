# ************************************************************
# Sequel Pro SQL 导出
# 版本 4468
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# 主机: 127.0.0.1 (MySQL 5.7.19)
# 数据库: godmin
# 生成时间: 2019-09-12 04:16:47 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 表 goadmin_menu 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_menu`;

CREATE TABLE `goadmin_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  -- 菜单ID，自增主键
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0',  -- 父级菜单ID，0表示顶级菜单
  `type` tinyint(4) unsigned NOT NULL DEFAULT '0',  -- 菜单类型
  `order` int(11) unsigned NOT NULL DEFAULT '0',  -- 菜单排序
  `title` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 菜单标题
  `icon` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 菜单图标
  `uri` varchar(3000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',  -- 菜单URI路径
  `header` varchar(150) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 菜单头部
  `plugin_name` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',  -- 插件名称
  `uuid` varchar(150) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 唯一标识符
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_menu` WRITE;
/*!40000 ALTER TABLE `goadmin_menu` DISABLE KEYS */;

INSERT INTO `goadmin_menu` (`id`, `parent_id`, `type`, `order`, `title`, `icon`, `uri`, `plugin_name`, `header`, `created_at`, `updated_at`)
VALUES
	(1,0,1,2,'Admin','fa-tasks','','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员菜单
	(2,1,1,2,'Users','fa-users','/info/manager','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 用户管理
	(3,1,1,3,'Roles','fa-user','/info/roles','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 角色管理
	(4,1,1,4,'Permission','fa-ban','/info/permission','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 权限管理
	(5,1,1,5,'Menu','fa-bars','/menu','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 菜单管理
	(6,1,1,6,'Operation log','fa-history','/info/op','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 操作日志
	(7,0,1,1,'Dashboard','fa-bar-chart','/','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00');  -- 仪表盘

/*!40000 ALTER TABLE `goadmin_menu` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_operation_log 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_operation_log`;

CREATE TABLE `goadmin_operation_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  -- 日志ID，自增主键
  `user_id` int(11) unsigned NOT NULL,  -- 用户ID
  `path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 请求路径
  `method` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 请求方法（GET/POST/PUT/DELETE等）
  `ip` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 客户端IP地址
  `input` text COLLATE utf8mb4_unicode_ci NOT NULL,  -- 请求输入参数
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`),
  KEY `admin_operation_log_user_id_index` (`user_id`)  -- 用户ID索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


# 表 goadmin_site 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_site`;

CREATE TABLE `goadmin_site` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,  -- 配置ID，自增主键
  `key` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 配置键名
  `value` longtext COLLATE utf8mb4_unicode_ci,  -- 配置值
  `description` varchar(3000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 配置描述
  `state` tinyint(3) unsigned NOT NULL DEFAULT '0',  -- 状态
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


# 表 goadmin_permissions 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_permissions`;

CREATE TABLE `goadmin_permissions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  -- 权限ID，自增主键
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 权限名称
  `slug` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 权限标识符
  `http_method` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- HTTP方法（GET/POST/PUT/DELETE等）
  `http_path` text COLLATE utf8mb4_unicode_ci NOT NULL,  -- HTTP路径
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_permissions_name_unique` (`name`)  -- 权限名称唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_permissions` DISABLE KEYS */;

INSERT INTO `goadmin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `created_at`, `updated_at`)
VALUES
	(1,'All permission','*','','*','2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 所有权限
	(2,'Dashboard','dashboard','GET,PUT,POST,DELETE','/','2019-09-10 00:00:00','2019-09-10 00:00:00');  -- 仪表盘权限

/*!40000 ALTER TABLE `goadmin_permissions` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_role_menu 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_role_menu`;

CREATE TABLE `goadmin_role_menu` (
  `role_id` int(11) unsigned NOT NULL,  -- 角色ID
  `menu_id` int(11) unsigned NOT NULL,  -- 菜单ID
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  KEY `admin_role_menu_role_id_menu_id_index` (`role_id`,`menu_id`)  -- 角色ID和菜单ID联合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_role_menu` WRITE;
/*!40000 ALTER TABLE `goadmin_role_menu` DISABLE KEYS */;

INSERT INTO `goadmin_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
	(1,1,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员角色拥有Admin菜单
	(1,7,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员角色拥有Dashboard菜单
	(2,7,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 操作员角色拥有Dashboard菜单
	(1,8,'2019-09-11 10:20:55','2019-09-11 10:20:55'),  -- 管理员角色拥有其他菜单
	(2,8,'2019-09-11 10:20:55','2019-09-11 10:20:55');  -- 操作员角色拥有其他菜单

/*!40000 ALTER TABLE `goadmin_role_menu` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_role_permissions 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_role_permissions`;

CREATE TABLE `goadmin_role_permissions` (
  `role_id` int(11) unsigned NOT NULL,  -- 角色ID
  `permission_id` int(11) unsigned NOT NULL,  -- 权限ID
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  UNIQUE KEY `admin_role_permissions` (`role_id`,`permission_id`)  -- 角色ID和权限ID唯一联合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_role_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_role_permissions` DISABLE KEYS */;

INSERT INTO `goadmin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
	(1,1,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员角色拥有所有权限
	(1,2,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员角色拥有仪表盘权限
	(2,2,'2019-09-10 00:00:00','2019-09-10 00:00:00');  -- 操作员角色拥有仪表盘权限

/*!40000 ALTER TABLE `goadmin_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_role_users 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_role_users`;

CREATE TABLE `goadmin_role_users` (
  `role_id` int(11) unsigned NOT NULL,  -- 角色ID
  `user_id` int(11) unsigned NOT NULL,  -- 用户ID
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  UNIQUE KEY `admin_user_roles` (`role_id`,`user_id`)  -- 角色ID和用户ID唯一联合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_role_users` WRITE;
/*!40000 ALTER TABLE `goadmin_role_users` DISABLE KEYS */;

INSERT INTO `goadmin_role_users` (`role_id`, `user_id`, `created_at`, `updated_at`)
VALUES
	(1,1,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- admin用户属于管理员角色
	(2,2,'2019-09-10 00:00:00','2019-09-10 00:00:00');  -- operator用户属于操作员角色

/*!40000 ALTER TABLE `goadmin_role_users` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_roles 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_roles`;

CREATE TABLE `goadmin_roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  -- 角色ID，自增主键
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 角色名称
  `slug` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 角色标识符
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_roles_name_unique` (`name`)  -- 角色名称唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_roles` WRITE;
/*!40000 ALTER TABLE `goadmin_roles` DISABLE KEYS */;

INSERT INTO `goadmin_roles` (`id`, `name`, `slug`, `created_at`, `updated_at`)
VALUES
	(1,'Administrator','administrator','2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员角色
	(2,'Operator','operator','2019-09-10 00:00:00','2019-09-10 00:00:00');  -- 操作员角色

/*!40000 ALTER TABLE `goadmin_roles` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_session 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_session`;

CREATE TABLE `goadmin_session` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,  -- 会话ID，自增主键
  `sid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',  -- 会话标识符
  `values` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',  -- 会话值
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# 表 goadmin_user_permissions 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_user_permissions`;

CREATE TABLE `goadmin_user_permissions` (
  `user_id` int(11) unsigned NOT NULL,  -- 用户ID
  `permission_id` int(11) unsigned NOT NULL,  -- 权限ID
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  UNIQUE KEY `admin_user_permissions` (`user_id`,`permission_id`)  -- 用户ID和权限ID唯一联合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_user_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_user_permissions` DISABLE KEYS */;

INSERT INTO `goadmin_user_permissions` (`user_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
	(1,1,'2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- admin用户拥有所有权限
	(2,2,'2019-09-10 00:00:00','2019-09-10 00:00:00');  -- operator用户拥有仪表盘权限

/*!40000 ALTER TABLE `goadmin_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;


# 表 goadmin_users 的数据导出
# ------------------------------------------------------------

DROP TABLE IF EXISTS `goadmin_users`;

CREATE TABLE `goadmin_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  -- 用户ID，自增主键
  `username` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 用户名
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',  -- 密码（加密）
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 显示名称
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 头像URL
  `remember_token` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 记住登录令牌
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_users_username_unique` (`username`)  -- 用户名唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `goadmin_users` WRITE;
/*!40000 ALTER TABLE `goadmin_users` DISABLE KEYS */;

INSERT INTO `goadmin_users` (`id`, `username`, `password`, `name`, `avatar`, `remember_token`, `created_at`, `updated_at`)
VALUES
	(1,'admin','$2a$10$U3F/NSaf2kaVbyXTBp7ppOn0jZFyRqXRnYXB.AMioCjXl3Ciaj4oy','admin','','tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh','2019-09-10 00:00:00','2019-09-10 00:00:00'),  -- 管理员账户（默认密码：admin）
	(2,'operator','$2a$10$rVqkOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.','Operator','',NULL,'2019-09-10 00:00:00','2019-09-10 00:00:00');  -- 操作员账户

/*!40000 ALTER TABLE `goadmin_users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
