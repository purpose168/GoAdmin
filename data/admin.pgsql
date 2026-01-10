--
-- PostgreSQL 数据库导出
--

-- 导出数据库版本 9.5.14
-- 导出工具 pg_dump 版本 10.5

SET statement_timeout = 0;  -- 设置语句超时时间为0（无限制）
SET lock_timeout = 0;  -- 设置锁超时时间为0（无限制）
SET idle_in_transaction_session_timeout = 0;  -- 设置事务空闲会话超时时间为0（无限制）
SET client_encoding = 'UTF8';  -- 设置客户端编码为UTF8
SET standard_conforming_strings = on;  -- 启用标准符合字符串
SELECT pg_catalog.set_config('search_path', '', false);  -- 设置搜索路径为空
SET check_function_bodies = false;  -- 禁用函数体检查
SET client_min_messages = warning;  -- 设置客户端最小消息级别为警告
SET row_security = off;  -- 禁用行级安全

--
-- 名称: plpgsql; 类型: EXTENSION; 模式: -; 所有者:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;  -- 如果不存在则创建 plpgsql 扩展


--
-- 名称: EXTENSION plpgsql; 类型: COMMENT; 模式: -; 所有者:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL 过程语言';  -- 为 plpgsql 扩展添加注释


--
-- 名称: goadmin_menu_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_menu_myid_seq  -- 创建菜单表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_menu_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

SET default_tablespace = '';  -- 设置默认表空间为空

SET default_with_oids = false;  -- 禁用 OID（对象标识符）

--
-- 名称: goadmin_menu; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_menu (  -- 创建菜单表
    id integer DEFAULT nextval('public.goadmin_menu_myid_seq'::regclass) NOT NULL,  -- 菜单ID，使用序列自增
    parent_id integer DEFAULT 0 NOT NULL,  -- 父级菜单ID，0表示顶级菜单
    type integer DEFAULT 0,  -- 菜单类型
    "order" integer DEFAULT 0 NOT NULL,  -- 菜单排序
    title character varying(50) NOT NULL,  -- 菜单标题
    header character varying(100),  -- 菜单头部
    plugin_name character varying(100) NOT NULL,  -- 插件名称
    icon character varying(50) NOT NULL,  -- 菜单图标
    uri character varying(3000) NOT NULL,  -- 菜单URI路径
    uuid character varying(100),  -- 唯一标识符
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_menu OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_operation_log_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_operation_log_myid_seq  -- 创建操作日志表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_operation_log_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_operation_log; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_operation_log (  -- 创建操作日志表
    id integer DEFAULT nextval('public.goadmin_operation_log_myid_seq'::regclass) NOT NULL,  -- 日志ID，使用序列自增
    user_id integer NOT NULL,  -- 用户ID
    path character varying(255) NOT NULL,  -- 请求路径
    method character varying(10) NOT NULL,  -- 请求方法（GET/POST/PUT/DELETE等）
    ip character varying(15) NOT NULL,  -- 客户端IP地址
    input text NOT NULL,  -- 请求输入参数
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_operation_log OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_site_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_site_myid_seq  -- 创建站点配置表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_site_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_site; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_site (  -- 创建站点配置表
    id integer DEFAULT nextval('public.goadmin_site_myid_seq'::regclass) NOT NULL,  -- 配置ID，使用序列自增
    key character varying(100) NOT NULL,  -- 配置键名
    value text NOT NULL,  -- 配置值
    type integer DEFAULT 0,  -- 配置类型
    description character varying(3000),  -- 配置描述
    state integer DEFAULT 0,  -- 状态
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_site OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_permissions_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_permissions_myid_seq  -- 创建权限表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_permissions_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_permissions; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_permissions (  -- 创建权限表
    id integer DEFAULT nextval('public.goadmin_permissions_myid_seq'::regclass) NOT NULL,  -- 权限ID，使用序列自增
    name character varying(50) NOT NULL,  -- 权限名称
    slug character varying(50) NOT NULL,  -- 权限标识符
    http_method character varying(255),  -- HTTP方法（GET/POST/PUT/DELETE等）
    http_path text NOT NULL,  -- HTTP路径
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_permissions OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_role_menu; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_role_menu (  -- 创建角色菜单关联表
    role_id integer NOT NULL,  -- 角色ID
    menu_id integer NOT NULL,  -- 菜单ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_role_menu OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_role_permissions; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_role_permissions (  -- 创建角色权限关联表
    role_id integer NOT NULL,  -- 角色ID
    permission_id integer NOT NULL,  -- 权限ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_role_permissions OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_role_users; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_role_users (  -- 创建角色用户关联表
    role_id integer NOT NULL,  -- 角色ID
    user_id integer NOT NULL,  -- 用户ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_role_users OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_roles_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_roles_myid_seq  -- 创建角色表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_roles_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_roles; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_roles (  -- 创建角色表
    id integer DEFAULT nextval('public.goadmin_roles_myid_seq'::regclass) NOT NULL,  -- 角色ID，使用序列自增
    name character varying NOT NULL,  -- 角色名称
    slug character varying NOT NULL,  -- 角色标识符
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_roles OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_session_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_session_myid_seq  -- 创建会话表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_session_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_session; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_session (  -- 创建会话表
    id integer DEFAULT nextval('public.goadmin_session_myid_seq'::regclass) NOT NULL,  -- 会话ID，使用序列自增
    sid character varying(50) NOT NULL,  -- 会话标识符
    "values" character varying(3000) NOT NULL,  -- 会话值
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_session OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_user_permissions; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_user_permissions (  -- 创建用户权限关联表
    user_id integer NOT NULL,  -- 用户ID
    permission_id integer NOT NULL,  -- 权限ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_user_permissions OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 名称: goadmin_users_myid_seq; 类型: SEQUENCE; 模式: public; 所有者: postgres
--

CREATE SEQUENCE public.goadmin_users_myid_seq  -- 创建用户表的自增序列
    START WITH 1  -- 起始值为1
    INCREMENT BY 1  -- 每次递增1
    NO MINVALUE  -- 无最小值
    MAXVALUE 99999999  -- 最大值为99999999
    CACHE 1;  -- 缓存1个序列值


ALTER TABLE public.goadmin_users_myid_seq OWNER TO postgres;  -- 将序列所有者设置为 postgres

--
-- 名称: goadmin_users; 类型: TABLE; 模式: public; 所有者: postgres
--

CREATE TABLE public.goadmin_users (  -- 创建用户表
    id integer DEFAULT nextval('public.goadmin_users_myid_seq'::regclass) NOT NULL,  -- 用户ID，使用序列自增
    username character varying(100) NOT NULL,  -- 用户名
    password character varying(100) NOT NULL,  -- 密码（加密）
    name character varying(100) NOT NULL,  -- 显示名称
    avatar character varying(255),  -- 头像URL
    remember_token character varying(100),  -- 记住登录令牌
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间，默认为当前时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间，默认为当前时间
);


ALTER TABLE public.goadmin_users OWNER TO postgres;  -- 将表所有者设置为 postgres

--
-- 数据: goadmin_menu; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_menu (id, parent_id, type, "order", title, plugin_name, header, icon, uri, created_at, updated_at) FROM stdin;  -- 导入菜单表数据
1	0	1	2	Admin		\N	fa-tasks		2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员菜单
2	1	1	2	Users		\N	fa-users	/info/manager	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 用户管理
3	1	1	3	Roles		\N	fa-user	/info/roles	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 角色管理
4	1	1	4	Permission		\N	fa-ban	/info/permission	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 权限管理
5	1	1	5	Menu		\N	fa-bars	/menu	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 菜单管理
6	1	1	6	Operation log		\N	fa-history	/info/op	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 操作日志
7	0	1	1	Dashboard		\N	fa-bar-chart	/	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 仪表盘
\.


--
-- 数据: goadmin_operation_log; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_operation_log (id, user_id, path, method, ip, input, created_at, updated_at) FROM stdin;  -- 导入操作日志表数据（无数据）
\.


--
-- 数据: goadmin_site; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_site (id, key, value, description, state, created_at, updated_at) FROM stdin;  -- 导入站点配置表数据（无数据）
\.


--
-- 数据: goadmin_permissions; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_permissions (id, name, slug, http_method, http_path, created_at, updated_at) FROM stdin;  -- 导入权限表数据
1	All permission	*		*	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 所有权限
2	Dashboard	dashboard	GET,PUT,POST,DELETE	/	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 仪表盘权限
\.


--
-- 数据: goadmin_role_menu; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_role_menu (role_id, menu_id, created_at, updated_at) FROM stdin;  -- 导入角色菜单关联表数据
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员角色拥有Admin菜单
1	7	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员角色拥有Dashboard菜单
2	7	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 操作员角色拥有Dashboard菜单
\.


--
-- 数据: goadmin_role_permissions; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_role_permissions (role_id, permission_id, created_at, updated_at) FROM stdin;  -- 导入角色权限关联表数据
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员角色拥有所有权限
1	2	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员角色拥有仪表盘权限
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 操作员角色拥有仪表盘权限
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
0	3	\N	\N  -- 测试数据（可忽略）
\.


--
-- 数据: goadmin_role_users; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_role_users (role_id, user_id, created_at, updated_at) FROM stdin;  -- 导入角色用户关联表数据
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00  -- admin用户属于管理员角色
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00  -- operator用户属于操作员角色
\.


--
-- 数据: goadmin_roles; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_roles (id, name, slug, created_at, updated_at) FROM stdin;  -- 导入角色表数据
1	Administrator	administrator	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员角色
2	Operator	operator	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 操作员角色
\.


--
-- 数据: goadmin_session; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_session (id, sid, "values", created_at, updated_at) FROM stdin;  -- 导入会话表数据（无数据）
\.


--
-- 数据: goadmin_user_permissions; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_user_permissions (user_id, permission_id, created_at, updated_at) FROM stdin;  -- 导入用户权限关联表数据
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00  -- admin用户拥有所有权限
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00  -- operator用户拥有仪表盘权限
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
0	1	\N	\N  -- 测试数据（可忽略）
\.


--
-- 数据: goadmin_users; 类型: TABLE DATA; 模式: public; 所有者: postgres
--

COPY public.goadmin_users (id, username, password, name, avatar, remember_token, created_at, updated_at) FROM stdin;  -- 导入用户表数据
1	admin	$2a$10$OxWYJJGTP2gi00l2x06QuOWqw5VR47MQCJ0vNKnbMYfrutij10Hwe	admin		tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 管理员账户（默认密码：admin）
2	operator	$2a$10$rVqkOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.	Operator		\N	2019-09-10 00:00:00	2019-09-10 00:00:00  -- 操作员账户
\.


--
-- 名称: goadmin_menu_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_menu_myid_seq', 7, true);  -- 设置菜单序列的当前值为7


--
-- 名称: goadmin_operation_log_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_operation_log_myid_seq', 1, true);  -- 设置操作日志序列的当前值为1


--
-- 名称: goadmin_permissions_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_permissions_myid_seq', 2, true);  -- 设置权限序列的当前值为2


--
-- 名称: goadmin_roles_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_roles_myid_seq', 2, true);  -- 设置角色序列的当前值为2


--
-- 名称: goadmin_site_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_site_myid_seq', 1, true);  -- 设置站点配置序列的当前值为1


--
-- 名称: goadmin_session_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_session_myid_seq', 1, true);  -- 设置会话序列的当前值为1


--
-- 名称: goadmin_users_myid_seq; 类型: SEQUENCE SET; 模式: public; 所有者: postgres
--

SELECT pg_catalog.setval('public.goadmin_users_myid_seq', 2, true);  -- 设置用户序列的当前值为2


--
-- 名称: goadmin_menu goadmin_menu_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_menu
    ADD CONSTRAINT goadmin_menu_pkey PRIMARY KEY (id);  -- 为菜单表添加主键约束


--
-- 名称: goadmin_operation_log goadmin_operation_log_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_operation_log
    ADD CONSTRAINT goadmin_operation_log_pkey PRIMARY KEY (id);  -- 为操作日志表添加主键约束


--
-- 名称: goadmin_permissions goadmin_permissions_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_permissions
    ADD CONSTRAINT goadmin_permissions_pkey PRIMARY KEY (id);  -- 为权限表添加主键约束


--
-- 名称: goadmin_roles goadmin_roles_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_roles
    ADD CONSTRAINT goadmin_roles_pkey PRIMARY KEY (id);  -- 为角色表添加主键约束


--
-- 名称: goadmin_site goadmin_site_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_site
    ADD CONSTRAINT goadmin_site_pkey PRIMARY KEY (id);  -- 为站点配置表添加主键约束


--
-- 名称: goadmin_session goadmin_session_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_session
    ADD CONSTRAINT goadmin_session_pkey PRIMARY KEY (id);  -- 为会话表添加主键约束


--
-- 名称: goadmin_users goadmin_users_pkey; 类型: CONSTRAINT; 模式: public; 所有者: postgres
--

ALTER TABLE ONLY public.goadmin_users
    ADD CONSTRAINT goadmin_users_pkey PRIMARY KEY (id);  -- 为用户表添加主键约束


--
-- 名称: SCHEMA public; 类型: ACL; 模式: -; 所有者: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;  -- 撤销 public 模式的所有权限
REVOKE ALL ON SCHEMA public FROM postgres;  -- 撤销 postgres 用户对 public 模式的所有权限
GRANT ALL ON SCHEMA public TO postgres;  -- 授予 postgres 用户对 public 模式的所有权限
GRANT ALL ON SCHEMA public TO PUBLIC;  -- 授予所有用户对 public 模式的所有权限


--
-- PostgreSQL 数据库导出完成
--
