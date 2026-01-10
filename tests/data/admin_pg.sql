-- PostgreSQL数据库转储
-- PostgreSQL database dump

-- 从数据库版本 9.5.14 转储
-- Dumped from database version 9.5.14
-- 由 pg_dump 版本 10.5 转储
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'EUC_CN';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

-- 名称: plpgsql; 类型: 扩展; 模式: -; 所有者: 
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


-- 名称: EXTENSION plpgsql; 类型: 注释; 模式: -; 所有者: 
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL过程语言';
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


-- 名称: goadmin_menu_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_menu_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_menu_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_menu_myid_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

-- 名称: goadmin_menu; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_menu; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_menu (
    id integer DEFAULT nextval('public.goadmin_menu_myid_seq'::regclass) NOT NULL,  -- 主键ID
    parent_id integer DEFAULT 0 NOT NULL,  -- 父级菜单ID
    type integer DEFAULT 0,  -- 菜单类型
    "order" integer DEFAULT 0 NOT NULL,  -- 排序
    title character varying(50) NOT NULL,  -- 菜单标题
    header character varying(100),  -- 菜单头部
    icon character varying(50) NOT NULL,  -- 图标
    uri character varying(50) NOT NULL,  -- 路径
    uuid character varying(100),  -- UUID
    plugin_name character varying(150) NOT NULL,  -- 插件名称
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_menu OWNER TO postgres;

-- 名称: goadmin_operation_log_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_operation_log_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_operation_log_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_operation_log_myid_seq OWNER TO postgres;

-- 名称: goadmin_operation_log; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_operation_log; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_operation_log (
    id integer DEFAULT nextval('public.goadmin_operation_log_myid_seq'::regclass) NOT NULL,  -- 主键ID
    user_id integer NOT NULL,  -- 用户ID
    path character varying(255) NOT NULL,  -- 请求路径
    method character varying(10) NOT NULL,  -- 请求方法
    ip character varying(15) NOT NULL,  -- IP地址
    input text NOT NULL,  -- 输入参数
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_operation_log OWNER TO postgres;

-- 名称: goadmin_permissions_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_permissions_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_permissions_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_permissions_myid_seq OWNER TO postgres;

-- 名称: goadmin_permissions; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_permissions; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_permissions (
    id integer DEFAULT nextval('public.goadmin_permissions_myid_seq'::regclass) NOT NULL,  -- 主键ID
    name character varying(50) NOT NULL,  -- 权限名称
    slug character varying(50) NOT NULL,  -- 权限标识
    http_method character varying(255),  -- HTTP方法
    http_path text NOT NULL,  -- HTTP路径
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_permissions OWNER TO postgres;

-- 名称: goadmin_site_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_site_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_site_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_site_myid_seq OWNER TO postgres;

-- 名称: goadmin_site; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_site; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_site (
    id integer DEFAULT nextval('public.goadmin_site_myid_seq'::regclass) NOT NULL,  -- 主键ID
    key character varying(100) NOT NULL,  -- 键
    value text NOT NULL,  -- 值
    type integer DEFAULT 0,  -- 类型
    description character varying(3000),  -- 描述
    state integer DEFAULT 0,  -- 状态
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_site OWNER TO postgres;

-- 名称: goadmin_role_menu; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_role_menu; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_role_menu (
    role_id integer NOT NULL,  -- 角色ID
    menu_id integer NOT NULL,  -- 菜单ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_role_menu OWNER TO postgres;

-- 名称: goadmin_role_permissions; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_role_permissions; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_role_permissions (
    role_id integer NOT NULL,  -- 角色ID
    permission_id integer NOT NULL,  -- 权限ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_role_permissions OWNER TO postgres;

-- 名称: goadmin_role_users; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_role_users; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_role_users (
    role_id integer NOT NULL,  -- 角色ID
    user_id integer NOT NULL,  -- 用户ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_role_users OWNER TO postgres;

-- 名称: goadmin_roles_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_roles_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_roles_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_roles_myid_seq OWNER TO postgres;

-- 名称: goadmin_roles; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_roles; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_roles (
    id integer DEFAULT nextval('public.goadmin_roles_myid_seq'::regclass) NOT NULL,  -- 主键ID
    name character varying NOT NULL,  -- 角色名称
    slug character varying NOT NULL,  -- 角色标识
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_roles OWNER TO postgres;

-- 名称: goadmin_session_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_session_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_session_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_session_myid_seq OWNER TO postgres;

-- 名称: goadmin_session; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_session; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_session (
    id integer DEFAULT nextval('public.goadmin_session_myid_seq'::regclass) NOT NULL,  -- 主键ID
    sid character varying(50) NOT NULL,  -- 会话ID
    "values" character varying(3000) NOT NULL,  -- 会话值
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_session OWNER TO postgres;

-- 名称: goadmin_user_permissions; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_user_permissions; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_user_permissions (
    user_id integer NOT NULL,  -- 用户ID
    permission_id integer NOT NULL,  -- 权限ID
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_user_permissions OWNER TO postgres;

-- 名称: goadmin_users_myid_seq; 类型: 序列; 模式: public; 所有者: postgres
-- Name: goadmin_users_myid_seq; Type: SEQUENCE; Schema: public; Owner: postgres

CREATE SEQUENCE public.goadmin_users_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;


ALTER TABLE public.goadmin_users_myid_seq OWNER TO postgres;

-- 名称: goadmin_users; 类型: 表; 模式: public; 所有者: postgres
-- Name: goadmin_users; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.goadmin_users (
    id integer DEFAULT nextval('public.goadmin_users_myid_seq'::regclass) NOT NULL,  -- 主键ID
    username character varying(190) NOT NULL,  -- 用户名
    password character varying(80) NOT NULL,  -- 密码
    name character varying(255) NOT NULL,  -- 姓名
    avatar character varying(255),  -- 头像
    remember_token character varying(100),  -- 记住令牌
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.goadmin_users OWNER TO postgres;

-- 名称: user_like_books; 类型: 表; 模式: public; 所有者: postgres
-- Name: user_like_books; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.user_like_books (
    id integer,  -- 主键ID
    user_id integer,  -- 用户ID
    name character varying,  -- 名称
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now()  -- 更新时间
);


ALTER TABLE public.user_like_books OWNER TO postgres;

-- 名称: users; 类型: 表; 模式: public; 所有者: postgres
-- Name: users; Type: TABLE; Schema: public; Owner: postgres

CREATE TABLE public.users (
    id integer NOT NULL,  -- 主键ID
    name character varying(100),  -- 姓名
    homepage character varying(3000),  -- 主页
    email character varying(100),  -- 邮箱
    birthday timestamp with time zone,  -- 生日
    country character varying(50),  -- 国家
    city character varying(50),  -- 城市
    password character varying(100),  -- 密码
    ip character varying(20),  -- IP地址
    certificate character varying(300),  -- 证书
    money integer,  -- 金额
    resume text,  -- 简历
    gender smallint,  -- 性别
    fruit character varying(200),  -- 水果
    drink character varying(200),  -- 饮料
    experience smallint,  -- 经验
    created_at timestamp without time zone DEFAULT now(),  -- 创建时间
    updated_at timestamp without time zone DEFAULT now(),  -- 更新时间
    member_id integer DEFAULT 0  -- 会员ID
);


ALTER TABLE public.users OWNER TO postgres;

-- goadmin_menu 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_menu; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_menu (id, parent_id, type, "order", title, plugin_name, header, icon, uri, created_at, updated_at) FROM stdin;
1	0	1	2	Admin		\N	fa-tasks		2019-09-10 00:00:00	2019-09-10 00:00:00
2	1	1	2	Users		\N	fa-users	/info/manager	2019-09-10 00:00:00	2019-09-10 00:00:00
3	1	1	3	Roles		\N	fa-user	/info/roles	2019-09-10 00:00:00	2019-09-10 00:00:00
4	1	1	4	Permission		\N	fa-ban	/info/permission	2019-09-10 00:00:00	2019-09-10 00:00:00
5	1	1	5	Menu		\N	fa-bars	/menu	2019-09-10 00:00:00	2019-09-10 00:00:00
6	1	1	6	Operation log		\N	fa-history	/info/op	2019-09-10 00:00:00	2019-09-10 00:00:00
7	0	1	1	Dashboard		\N	fa-bar-chart	/	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- goadmin_operation_log 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_operation_log; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_operation_log (id, user_id, path, method, ip, input, created_at, updated_at) FROM stdin;
\.


-- goadmin_site 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_site; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_site (id, key, value, description, state, created_at, updated_at) FROM stdin;
\.


-- goadmin_permissions 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_permissions; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_permissions (id, name, slug, http_method, http_path, created_at, updated_at) FROM stdin;
1	All permission	*		*	2019-09-10 00:00:00	2019-09-10 00:00:00
2	Dashboard	dashboard	GET,PUT,POST,DELETE	/	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- goadmin_role_menu 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_role_menu; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_role_menu (role_id, menu_id, created_at, updated_at) FROM stdin;
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00
1	7	2019-09-10 00:00:00	2019-09-10 00:00:00
2	7	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- goadmin_role_permissions 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_role_permissions (role_id, permission_id, created_at, updated_at) FROM stdin;
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00
1	2	2019-09-10 00:00:00	2019-09-10 00:00:00
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
0	3	\N	\N
\.


-- goadmin_role_users 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_role_users; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_role_users (role_id, user_id, created_at, updated_at) FROM stdin;
1	1	2019-09-10 00:00:00	2019-09-10 00:00:00
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- goadmin_roles 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_roles; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_roles (id, name, slug, created_at, updated_at) FROM stdin;
1	Administrator	administrator	2019-09-10 00:00:00	2019-09-10 00:00:00
2	Operator	operator	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- goadmin_session 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_session; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_session (id, sid, "values", created_at, updated_at) FROM stdin;
2	f5a99916-36c8-4fd6-8873-6f2be8845cd0	{"user_id":1}	2019-11-27 22:26:11.917665	2019-11-27 22:26:11.917665
3	03263ffc-0043-4b89-a02f-3aa616bbf857	{"user_id":3}	2019-11-27 22:26:12.819931	2019-11-27 22:26:12.819931
\.


-- goadmin_user_permissions 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_user_permissions; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_user_permissions (user_id, permission_id, created_at, updated_at) FROM stdin;
2	2	2019-09-10 00:00:00	2019-09-10 00:00:00
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
0	1	\N	\N
1	1	2019-11-27 22:26:12.425769	2019-11-27 22:26:12.425769
3	1	2019-11-27 22:26:12.572997	2019-11-27 22:26:12.572997
\.


-- goadmin_users 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: goadmin_users; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.goadmin_users (id, username, password, name, avatar, remember_token, created_at, updated_at) FROM stdin;
1	admin	$2a$10$OxWYJJGTP2gi00l2x06QuOWqw5VR47MQCJ0vNKnbMYfrutij10Hwe	admin		tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh	2019-09-10 00:00:00	2019-09-10 00:00:00
2	operator	$2a$10$rVqkOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.	Operator		\N	2019-09-10 00:00:00	2019-09-10 00:00:00
\.


-- user_like_books 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: user_like_books; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.user_like_books (id, user_id, name, created_at, updated_at) FROM stdin;
1	1	Robinson Crusoe	2020-03-15 09:00:57.409596	2020-03-15 09:00:57.409596
\.


-- users 的数据; 类型: 表数据; 模式: public; 所有者: postgres
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres

COPY public.users (id, name, homepage, email, birthday, country, city, password, ip, certificate, money, resume, gender, fruit, drink, experience, created_at, updated_at, member_id) FROM stdin;
1	Jack	http://jack.me	jack@163.com	1993-10-21 00:00:00+08	china	guangzhou	123456	127.0.0.1	\N	10	<h1>Jacks Resume</h1>	0	apple	water	0	2020-03-09 15:24:00	2020-03-09 15:24:00	0
\.


-- 名称: goadmin_menu_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_menu_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_menu_myid_seq', 7, true);


-- 名称: goadmin_operation_log_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_operation_log_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_operation_log_myid_seq', 1, true);


-- 名称: goadmin_permissions_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_permissions_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_permissions_myid_seq', 2, true);


-- 名称: goadmin_roles_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_roles_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_roles_myid_seq', 2, true);


-- 名称: goadmin_session_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_session_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_session_myid_seq', 1, true);

-- 名称: goadmin_site_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_site_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_site_myid_seq', 1, true);


-- 名称: goadmin_users_myid_seq; 类型: 序列设置; 模式: public; 所有者: postgres
-- Name: goadmin_users_myid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

SELECT pg_catalog.setval('public.goadmin_users_myid_seq', 2, true);


-- 名称: goadmin_menu goadmin_menu_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_menu goadmin_menu_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_menu
    ADD CONSTRAINT goadmin_menu_pkey PRIMARY KEY (id);


-- 名称: goadmin_operation_log goadmin_operation_log_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_operation_log goadmin_operation_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_operation_log
    ADD CONSTRAINT goadmin_operation_log_pkey PRIMARY KEY (id);


-- 名称: goadmin_permissions goadmin_permissions_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_permissions goadmin_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_permissions
    ADD CONSTRAINT goadmin_permissions_pkey PRIMARY KEY (id);


-- 名称: goadmin_roles goadmin_roles_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_roles goadmin_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_roles
    ADD CONSTRAINT goadmin_roles_pkey PRIMARY KEY (id);


-- 名称: goadmin_site goadmin_site_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_site goadmin_site_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_site
    ADD CONSTRAINT goadmin_site_pkey PRIMARY KEY (id);

-- 名称: goadmin_session goadmin_session_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_session goadmin_session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_session
    ADD CONSTRAINT goadmin_session_pkey PRIMARY KEY (id);


-- 名称: goadmin_users goadmin_users_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: goadmin_users goadmin_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.goadmin_users
    ADD CONSTRAINT goadmin_users_pkey PRIMARY KEY (id);


-- 名称: users users_pkey; 类型: 约束; 模式: public; 所有者: postgres
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


-- 名称: SCHEMA public; 类型: ACL; 模式: -; 所有者: postgres
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- PostgreSQL数据库转储完成
-- PostgreSQL database dump complete

