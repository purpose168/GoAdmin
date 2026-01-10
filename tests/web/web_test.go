// web_test.go - Web UI 自动化测试文件
// 包名：web
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含 GoAdmin 管理后台的 Web UI 自动化测试用例，使用 Selenium 进行浏览器自动化测试
//       涵盖登录、数据表格操作、表单操作、权限管理、菜单管理、用户管理等核心功能测试

package web

import (
	"io"
	"log"
	"os"
	"testing"

	_ "github.com/purpose168/GoAdmin-themes/adminlte"          // 导入 AdminLTE 主题
	_ "github.com/purpose168/GoAdmin/adapter/gin"              // 导入 Gin 适配器
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql" // 导入 MySQL 数据库驱动

	"github.com/gin-gonic/gin"                       // Gin Web 框架
	"github.com/purpose168/GoAdmin/engine"           // GoAdmin 核心引擎
	"github.com/purpose168/GoAdmin/modules/config"   // 配置模块
	"github.com/purpose168/GoAdmin/plugins/admin"    // 管理插件
	"github.com/purpose168/GoAdmin/template"         // 模板引擎
	"github.com/purpose168/GoAdmin/template/chartjs" // Chart.js 图表组件
	"github.com/purpose168/GoAdmin/tests/tables"     // 测试数据表
)

// ==================== 页面元素 XPath 选择器常量定义 ====================
// 以下常量定义了页面中各个元素的 XPath 路径，用于 Selenium 自动化测试定位元素

const (
	// ==================== 信息表格页面元素 ====================

	newPageBtn              = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[3]/a`                                                            // 新建页面按钮
	editPageBtn             = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[8]/div/ul/li[1]/a`                                  // 编辑页面按钮
	genderActionDropDown    = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[7]/div/span/span[1]/span/span[2]`                                // 性别筛选下拉菜单
	menOptionActionBtn      = `/html/body/span/span/span[2]/ul/li[2]`                                                                                       // 男性选项按钮
	idOrderBtn              = `//*[@id="sort-id"]`                                                                                                          // ID 排序按钮
	rowActionDropDown       = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[8]/div/div/a`                                       // 行操作下拉菜单
	popupBtn                = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[5]/a`                                                            // 弹窗按钮
	popup                   = `//*[@id="pjax-container"]/section[2]/div/div/div[4]/div[3]`                                                                  // 弹窗容器
	popupCloseBtn           = `//*[@id="pjax-container"]/section[2]/div/div/div[4]/div[3]/div/div/div[3]/button`                                            // 弹窗关闭按钮
	ajaxBtn                 = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[6]/a`                                                            // Ajax 请求按钮
	ajaxAlert               = `/html/body/div[3]`                                                                                                           // Ajax 提示框
	selectionDropDown       = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[2]/div/div/div[1]/div/span/span[1]/span/span[2]`     // 单选下拉菜单
	selectionLi1            = `/html/body/span/span/span[2]/ul/li[1]`                                                                                       // 下拉选项 1
	selectionLi2            = `/html/body/span/span/span[2]/ul/li[2]`                                                                                       // 下拉选项 2
	selectionRes            = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[2]/div/div/div[1]/div/span/span[1]/span/span[1]`     // 单选结果显示
	multiSelectInput        = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[1]/div/div/div[2]/div/span/span[1]/span/ul/li/input` // 多选输入框
	multiSelectLi1          = `/html/body/span/span/span/ul/li[1]`                                                                                          // 多选选项 1
	multiSelectLi2          = `/html/body/span/span/span/ul/li[2]`                                                                                          // 多选选项 2
	multiSelectLi3          = `/html/body/span/span/span/ul/li[3]`                                                                                          // 多选选项 3
	multiSelectRes          = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[1]/div/div/div[2]/div/span/span[1]/span/ul/li[1]`    // 多选结果显示
	filterNameField         = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[1]/div/div/div[1]/div/div/input`                     // 名称筛选字段
	filterCreatedStart      = `//*[@id="created_at_start__goadmin"]`                                                                                        // 创建时间起始筛选
	filterCreatedEnd        = `//*[@id="created_at_end__goadmin"]`                                                                                          // 创建时间结束筛选
	radio                   = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div[3]/div/div/div[1]/div/div[1]`                        // 单选按钮
	searchBtn               = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[2]/div[2]/div[1]/button`                                        // 搜索按钮
	filterResetBtn          = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[2]/div[2]/div[2]/a`                                             // 筛选重置按钮
	rowSelector             = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[1]/button`                                                       // 行选择器
	rowSelectCityCheckbox   = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[1]/ul/li[1]/ul/li[4]/label/div`                                  // 城市复选框
	rowSelectAvatarCheckbox = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/div/div[1]/ul/li[1]/ul/li[5]/label/div`                                  // 头像复选框
	actionDropDown          = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[1]/div`                                             // 操作下拉菜单
	exportBtn               = `//*[@id="pjax-container"]/section[2]/div/div/div[1]/span/div/button`                                                         // 导出按钮
	previewAction           = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[8]/div/ul/li[7]/a`                                  // 预览操作
	closePreviewAction      = `//*[@id="pjax-container"]/section[2]/div/div/div[4]/div[2]/div/div/div[3]/button`                                            // 关闭预览操作
	previewPopup            = `//*[@id="pjax-container"]/section[2]/div/div/div[4]/div[2]`                                                                  // 预览弹窗
	rowAjaxAction           = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[8]/div/ul/li[6]/a`                                  // 行 Ajax 操作
	rowAjaxPopup            = `/html/body/div[3]`                                                                                                           // 行 Ajax 弹窗
	closeRowAjaxPopup       = `/html/body/div[3]/div[7]/div/button`                                                                                         // 关闭行 Ajax 弹窗
	updateNameTd            = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[3]/a`                                               // 更新名称单元格
	updateNameInput         = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[3]/div/div[2]/div/form/div/div[1]/div[1]/input`     // 更新名称输入框
	updateNameSaveBtn       = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[3]/div/div[2]/div/form/div/div[1]/div[2]/button[1]` // 更新名称保存按钮
	updateGenderBtn         = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[4]/div/div/span[1]`                                 // 更新性别按钮
	detailBtn               = `//*[@id="pjax-container"]/section[2]/div/div/div[3]/table/tbody/tr[2]/td[10]/a`                                              // 详情按钮

	// ==================== 表单页面元素 ====================

	saveBtn            = `//*[@id="pjax-container"]/section[2]/div/form/div[2]/div[2]/div[1]/button` // 保存按钮
	resetBtn           = `//*[@id="pjax-container"]/section[2]/div/form/div[2]/div[2]/div[2]/button` // 重置按钮
	nameField          = `//*[@id="tab-form-0"]/div[1]/div/div/input`                                // 姓名字段
	ageField           = `//*[@id="tab-form-0"]/div[2]/div/div/div/input`                            // 年龄字段
	emailField         = `//*[@id="tab-form-0"]/div[4]/div/div/input`                                // 邮箱字段
	birthdayField      = `//*[@id="tab-form-0"]/div[5]/div/div/input`                                // 生日字段
	passwordField      = `//*[@id="tab-form-0"]/div[6]/div/div/input`                                // 密码字段
	homePageField      = `//*[@id="tab-form-0"]/div[3]/div/div/input`                                // 主页字段
	ipField            = `//*[@id="tab-form-0"]/div[7]/div/div/input`                                // IP 地址字段
	amountField        = `//*[@id="tab-form-0"]/div[9]/div/div/input`                                // 金额字段
	appleOptField      = `//*[@id="bootstrap-duallistbox-nonselected-list_fruit[]"]/option[1]`       // 苹果选项
	bananaOptField     = `//*[@id="bootstrap-duallistbox-nonselected-list_fruit[]"]/option[2]`       // 香蕉选项
	watermelonOptField = `//*[@id="bootstrap-duallistbox-nonselected-list_fruit[]"]/option[3]`       // 西瓜选项
	//pearOptField          = `//*[@id="bootstrap-duallistbox-nonselected-list_fruit[]"]/option[4]` // 梨选项（已注释）
	genderBoyCheckBox     = `//*[@id="tab-form-1"]/div[5]/div/div/div[1]`                           // 男性复选框
	genderGirlCheckBox    = `//*[@id="tab-form-1"]/div[5]/div/div/div[2]`                           // 女性复选框
	experienceDropDown    = `//*[@id="tab-form-1"]/div[7]/div/span/span[1]/span/span[2]`            // 经验下拉菜单
	twoYearsSelection     = `/html/body/span/span/span[2]/ul/li[1]`                                 // 两年经验选项
	threeYearsSelection   = `/html/body/span/span/span[2]/ul/li[2]`                                 // 三年经验选项
	fourYearsSelection    = `/html/body/span/span/span[2]/ul/li[3]`                                 // 四年经验选项
	fiveYearsSelection    = `/html/body/span/span/span[2]/ul/li[4]`                                 // 五年经验选项
	inputTab              = `//*[@id="pjax-container"]/section[2]/div/form/div[1]/div/div/ul/li[1]` // 输入标签页
	selectTab             = `//*[@id="pjax-container"]/section[2]/div/form/div[1]/div/div/ul/li[2]` // 选择标签页
	multiSelectionInput   = `//*[@id="tab-form-1"]/div[6]/div/span/span[1]/span/ul/li[2]/input`     // 多选输入框
	multiSelectedOpt      = `//*[@id="tab-form-1"]/div[6]/div/span/span[1]/span/ul/li[1]`           // 多选已选项
	multiBeerOpt          = `/html/body/span/span/span/ul/li[1]`                                    // 啤酒选项
	multiJuiceOpt         = `/html/body/span/span/span/ul/li[2]`                                    // 果汁选项
	multiWaterOpt         = `/html/body/span/span/span/ul/li[3]`                                    // 水选项
	multiRedBullOpt       = `/html/body/span/span/span/ul/li[4]`                                    // 红牛选项
	continueEditCheckBox  = `//*[@id="pjax-container"]/section[2]/div/form/div[2]/div[2]/label/div` // 继续编辑复选框
	boxSelectedOpt        = `//*[@id="bootstrap-duallistbox-selected-list_fruit[]"]/option`         // 已选框选项
	experienceSelectedOpt = `//*[@id="tab-form-1"]/div[7]/div/span/span[1]/span/span[1]`            // 经验已选项

	sideBarManageDropDown    = `/html/body/div[1]/aside/section/ul/li[2]/a/span[2]`                                                                       // 侧边栏管理下拉菜单
	menuPageBtn              = `/html/body/div[1]/aside/section/ul/li[2]/ul/li[4]/a`                                                                      // 菜单页面按钮
	menuParentIdDropDown     = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[1]/div/span/span[1]/span/span[2]`  // 菜单父级 ID 下拉菜单
	parentIdRootOpt          = `/html/body/span/span/span[2]/ul/li[1]`                                                                                    // 根节点选项
	parentIdDashboardOpt     = `/html/body/span/span/span[2]/ul/li[2]`                                                                                    // 仪表板选项
	parentIdAdminOpt         = `/html/body/span/span/span[2]/ul/li[3]`                                                                                    // 管理员选项
	parentIdUserOpt          = `/html/body/span/span/span[2]/ul/li[4]`                                                                                    // 用户选项
	menuRoleDropDown         = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[6]/div/span/span[1]/span`          // 菜单角色下拉菜单
	menuRoleAdminOpt         = `/html/body/span/span/span/ul/li[1]`                                                                                       // 管理员角色选项
	menuRoleOperatorOpt      = `/html/body/span/span/span/ul/li[2]`                                                                                       // 操作员角色选项
	iconPopupBtn             = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[4]/div/div[1]/span`                // 图标弹窗按钮
	iconPopup                = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[4]/div/div[2]`                     // 图标弹窗
	iconBtn                  = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[4]/div/div[2]/div[3]/div/div/a[5]` // 图标按钮
	menuNameInput            = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[2]/div/div/input`                  // 菜单名称输入框
	menuUriInput             = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[1]/div/div/div[5]/div/div/input`                  // 菜单 URI 输入框
	menuInfoSaveBtn          = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/div/div[2]/form/div[2]/div[2]/div[1]/button`                          // 菜单信息保存按钮
	testMenuItem             = `//*[@id="tree-model"]/ol/li[2]`                                                                                           // 测试菜单项
	testMenuDeleteBtn        = `//*[@id="tree-model"]/ol/li[2]/div/span/a[2]`                                                                             // 测试菜单删除按钮
	testMenuDeleteConfirmBtn = `/html/body/div[3]/div[7]/div/button`                                                                                      // 测试菜单删除确认按钮
	menuOkBtn                = `/html/body/div[3]/div[7]/div/button`                                                                                      // 菜单确认按钮
	userMenuEditBtn          = `//*[@id="tree-model"]/ol/li[3]/div/span/a[1]`                                                                             // 用户菜单编辑按钮
	headFieldInput           = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[4]/div/div/input`                             // 头部字段输入框
	menuEditSaveBtn          = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[2]/div[2]/div[1]/button`                                     // 菜单编辑保存按钮

	managerPageBtn               = `/html/body/div[1]/aside/section/ul/li[2]/ul/li[1]/a`                                                              // 管理员页面按钮
	rolesPageBtn                 = `/html/body/div[1]/aside/section/ul/li[2]/ul/li[2]/a`                                                              // 角色页面按钮
	permissionPageBtn            = `/html/body/div[1]/aside/section/ul/li[2]/ul/li[3]/a`                                                              // 权限页面按钮
	operationLogPageBtn          = `/html/body/div[1]/aside/section/ul/li[2]/ul/li[5]/a`                                                              // 操作日志页面按钮
	navLinkBtn                   = `//*[@id="firstnav"]/div[2]/ul/li[1]/a`                                                                            // 导航链接按钮
	navCloseBtn                  = `//*[@id="firstnav"]/div[2]/ul/li[1]/i`                                                                            // 导航关闭按钮
	userPageBtn                  = `/html/body/div[1]/aside/section/ul/li[3]/a`                                                                       // 用户页面按钮
	managerEditBtn               = `//*[@id="pjax-container"]/section[2]/div/div[2]/div[2]/table/tbody/tr[3]/td[8]/a[1]`                              // 管理员编辑按钮
	operatorEditBtn              = `//*[@id="pjax-container"]/section[2]/div/div[2]/div[2]/table/tbody/tr[2]/td[8]/a[1]`                              // 操作员编辑按钮
	managerNameField             = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[2]/div/div/input`                     // 管理员名称字段
	managerNickNameField         = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[3]/div/div/input`                     // 管理员昵称字段
	managerRoleSelectedOpt       = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[5]/div/span[1]/span[1]/span/ul/li[1]` // 管理员角色已选项
	managerPermissionSelectedOpt = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[6]/div/span[1]/span[1]/span/ul/li[1]` // 管理员权限已选项
	managerRoleDropDown          = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[5]/div/span[1]/span[1]/span/ul`       // 管理员角色下拉菜单
	managerRoleOpt2              = `/html/body/span/span/span/ul/li[2]`                                                                               // 管理员角色选项 2
	managerPermissionDropDown    = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[6]/div/span[1]/span[1]/span`          // 管理员权限下拉菜单
	managerPermissionOpt2        = `/html/body/span/span/span/ul/li[2]`                                                                               // 管理员权限选项 2
	managerSaveBtn               = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[2]/div[2]/div[1]/button`                             // 管理员保存按钮
	newPermissionBtn             = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[6]/div/span[2]/a`                     // 新增权限按钮
	managerUserViewSelectOpt     = `/html/body/span/span/span/ul/li[3]`                                                                               // 管理员用户查看选择选项

	permissionNameInput    = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[1]/div/div/input`                        // 权限名称输入框
	permissionSlugInput    = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[2]/div/div/input`                        // 权限标识输入框
	permissionMethodSelect = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[3]/div/span[1]/span[1]/span/ul/li/input` // 权限方法选择
	permissionGetSelectOpt = `/html/body/span/span/span/ul/li[1]`                                                                                  // GET 方法选择选项
	permissionPathInput    = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[1]/div/div/div[4]/div/textarea`                         // 权限路径输入框
	permissionSaveBtn      = `//*[@id="pjax-container"]/section[2]/div/div/div[2]/form/div[2]/div[2]/div[1]/button`                                // 权限保存按钮

	userNavMenuBtn = `//*[@id="firstnav"]/div[4]/ul/li[5]/a`          // 用户导航菜单按钮
	userSettingBtn = `//*[@id="firstnav"]/div[4]/ul/li[5]/ul/li[5]/a` // 用户设置按钮
	userSignOutBtn = `//*[@id="firstnav"]/div[4]/ul/li[5]/ul/li[6]/a` // 用户退出按钮

	loginPageUserNameInput = `//*[@id="username"]` // 登录页面用户名输入框
	loginPagePasswordInput = `//*[@id="password"]` // 登录页面密码输入框
)

// ==================== 全局变量定义 ====================

var (
	debugMode  = false     // 调试模式开关，false 表示无头模式，true 表示显示浏览器窗口
	optionList = []string{ // Chrome 浏览器启动选项列表
		"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36", // 设置用户代理
		"--window-size=1500,900",              // 设置浏览器窗口大小
		"--incognito",                         // 使用无痕模式
		"--blink-settings=imagesEnabled=true", // 启用图片加载
		"--no-default-browser-check",          // 禁用默认浏览器检查
		"--ignore-ssl-errors=true",            // 忽略 SSL 错误
		"--ssl-protocol=any",                  // 允许任意 SSL 协议
		"--no-sandbox",                        // 禁用沙箱模式（Linux 下需要）
		"--disable-breakpad",                  // 禁用崩溃报告
		"--disable-gpu",                       // 禁用 GPU 加速
		"--disable-logging",                   // 禁用日志记录
		"--no-zygote",                         // 禁用 zygote 进程
		"--allow-running-insecure-content",    // 允许运行不安全内容
	}
)

// ==================== 常量定义 ====================

const (
	port = ":9033" // 服务器监听端口号
)

// ==================== 初始化函数 ====================
// init 函数在包加载时自动执行，用于初始化测试环境
func init() {
	// 检查命令行最后一个参数是否为 "true"，如果是则启用调试模式
	if os.Args[len(os.Args)-1] == "true" {
		debugMode = true
	}
	// 如果不是调试模式，添加 --headless 选项使浏览器在后台运行
	if !debugMode {
		optionList = append(optionList, "--headless")
	}
}

// ==================== 服务器启动函数 ====================
// startServer 启动测试用的 Web 服务器
// 参数：quit - 用于通知服务器退出的通道
// 说明：该函数启动一个 Gin Web 服务器，配置 GoAdmin 引擎和管理插件
func startServer(quit chan struct{}) {
	// 如果不是调试模式，设置 Gin 为发布模式并禁用日志输出
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	// 创建 Gin 路由实例
	r := gin.New()

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default()

	// 创建管理插件并添加数据表生成器
	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	// 添加 Chart.js 图表组件到模板引擎
	template.AddComp(chartjs.NewChart())

	// 从 JSON 配置文件读取配置
	cfg := config.ReadFromJson("./config.json")
	// 如果是调试模式，启用 SQL 日志、调试日志和访问日志
	if debugMode {
		cfg.SqlLog = true
		cfg.Debug = true
		cfg.AccessLogOff = false
	}

	// 配置引擎：添加配置、添加插件、使用 Gin 路由
	if err := eng.AddConfig(&cfg).
		AddPlugins(adminPlugin).
		Use(r); err != nil {
		panic(err)
	}

	// 注册 HTML 路由
	eng.HTML("GET", "/admin", tables.GetContent)

	// 配置静态文件服务
	r.Static("/uploads", "./uploads")

	// 在新的 goroutine 中启动 HTTP 服务器
	go func() {
		_ = r.Run(port)
	}()

	// 等待退出信号
	<-quit
	log.Print("closing database connection") // 打印关闭数据库连接日志
	eng.MysqlConnection().Close()            // 关闭 MySQL 数据库连接
}

// ==================== 主测试函数 ====================
// TestWeb 是 Web UI 自动化测试的主入口函数
// 参数：t - 测试对象
// 说明：该函数运行完整的用户验收测试套件，包括登录、表格操作、表单操作、权限管理等测试
func TestWeb(t *testing.T) {
	// 执行用户验收测试套件
	UserAcceptanceTestSuit(t, func(_ *testing.T, page *Page) {
		defer page.Destroy()
		testLogin(page)
		testInfoTablePageOperations(page)
		testNewPageOperations(page)
		testEditPageOperations(page)
		// testDetailPageOperations(page) // 详情页面操作测试（函数未定义，已注释）
		// testRolePageOperations(page) // 角色页面操作测试（函数未定义，已注释）
		// testPermissionPageOperations(page) // 权限页面操作测试（函数未定义，已注释）
		// testMenuPageOperations(page) // 菜单页面操作测试（函数未定义，已注释）
		// testManagerPageOperations(page) // 管理员页面操作测试（函数未定义，已注释）
		// testPermission(page) // 权限测试（函数未定义，已注释）
	}, startServer, debugMode, optionList...)
}

// ==================== 登录测试函数 ====================
// testLogin 测试用户登录功能
// 参数：page - 页面对象，用于操作浏览器
// 说明：该函数测试用户登录流程，包括填写用户名和密码、点击登录按钮、验证登录成功
func testLogin(page *Page) {
	page.NavigateTo(url("/login")) // 导航到登录页面

	page.Fill(loginPageUserNameInput, "admin") // 填写用户名
	page.Fill(loginPagePasswordInput, "admin") // 填写密码
	page.ClickS(page.FindByButton("Login"))    // 点击登录按钮

	wait(3) // 等待 3 秒

	page.Contain("main-header") // 验证页面包含 "main-header" 元素，表示登录成功
}

// ==================== 信息表格页面操作测试函数 ====================
// testInfoTablePageOperations 测试信息表格页面的各种操作
// 参数：page - 页面对象，用于操作浏览器
// 说明：该函数测试表格页面的按钮、筛选、排序、行选择、导出、操作按钮等功能
func testInfoTablePageOperations(page *Page) {
	// ==================== 导航链接检查 ====================
	// 以下代码用于检查侧边栏导航链接（已注释）

	//printPart("nav link check")
	//page.Click(sideBarManageDropDown)
	//page.Click(managerPageBtn)
	//page.Click(rolesPageBtn)
	//page.Click(permissionPageBtn)
	//page.Click(menuPageBtn)
	//page.Click(operationLogPageBtn)
	//page.Click(navLinkBtn)
	//page.Click(navCloseBtn)
	//page.Click(navLinkBtn)
	//page.Click(navCloseBtn)
	//page.Click(navLinkBtn)
	//page.Click(navCloseBtn)
	//page.Click(navLinkBtn)
	//page.Click(navCloseBtn)

	page.NavigateTo(url("/info/user")) // 导航到用户信息页面

	page.Contain("Users") // 验证页面包含 "Users" 文本

	// ==================== 按钮检查 ====================

	printPart("buttons check") // 打印测试部分信息

	page.Click(popupBtn) // 点击弹窗按钮

	page.Display(popup) // 验证弹窗显示

	wait(1) // 等待 1 秒

	page.Contain("hello world") // 验证弹窗包含 "hello world" 文本
	page.Click(popupCloseBtn)   // 点击关闭弹窗按钮

	page.Nondisplay(popup) // 验证弹窗已隐藏

	page.Click(ajaxBtn) // 点击 Ajax 按钮

	page.Contain("Oh li get") // 验证页面包含 "Oh li get" 文本

	page.Display(ajaxAlert) // 验证 Ajax 提示框显示

	page.ClickS(page.FindByButton("OK")) // 点击 OK 按钮

	page.Nondisplay(ajaxAlert) // 验证 Ajax 提示框已隐藏

	// ==================== 更新检查 ====================

	printPart("update check")                  // 打印测试部分信息
	page.Click(updateNameTd)                   // 点击名称单元格
	page.Fill(updateNameInput, "DukeDukeDuke") // 填写新名称
	page.Click(updateNameSaveBtn)              // 点击保存按钮
	page.Click(updateGenderBtn)                // 点击性别按钮
	page.Contain("DukeDukeDuke")               // 验证页面包含新名称

	// ==================== 筛选区域检查 ====================

	printPart("filter area check") // 打印测试部分信息

	page.Click(selectionDropDown) // 点击单选下拉菜单

	page.Text(selectionLi1, "men")   // 验证选项 1 文本为 "men"
	page.Text(selectionLi2, "women") // 验证选项 2 文本为 "women"

	page.Click(selectionLi2) // 点击 "women" 选项

	page.Attr(page.FindByXPath(selectionRes), "title", "women") // 验证结果显示的 title 属性为 "women"

	page.Fill(multiSelectInput, " ") // 在多选输入框中输入空格

	wait(1) // 等待 1 秒

	page.Text(multiSelectLi1, "water")    // 验证多选选项 1 文本为 "water"
	page.Text(multiSelectLi2, "juice")    // 验证多选选项 2 文本为 "juice"
	page.Text(multiSelectLi3, "red bull") // 验证多选选项 3 文本为 "red bull"

	page.Click(multiSelectLi3) // 点击 "red bull" 选项

	page.Attr(page.FindByXPath(multiSelectRes), "title", "red bull") // 验证结果显示的 title 属性为 "red bull"

	page.Click(radio) // 点击单选按钮

	page.Fill(filterNameField, "Jack") // 在名称筛选字段中填写 "Jack"

	//page.Fill(filterCreatedStart, "2020-03-08 15:24:00") // 填写创建时间起始（已注释）
	//page.Click(filterCreatedEnd) // 点击创建时间结束（已注释）

	page.Click(searchBtn, 2) // 点击搜索按钮，等待 2 秒

	page.Click(filterResetBtn, 2) // 点击筛选重置按钮，等待 2 秒

	// ==================== 行选择器检查 ====================

	printPart("row selector check") // 打印测试部分信息

	page.Click(rowSelector)             // 点击行选择器
	page.Click(rowSelectCityCheckbox)   // 点击城市复选框
	page.Click(rowSelectAvatarCheckbox) // 点击头像复选框

	page.ClickS(page.FindByButton("Submit"), 2) // 点击提交按钮，等待 2 秒

	page.NoContain("guangzhou") // 验证页面不包含 "guangzhou"

	page.ClickS(page.FindByID("filter-btn")) // 点击筛选按钮

	page.CssS(page.FindByClass("filter-area"), "display", "none") // 验证筛选区域的 display 样式为 "none"

	// ==================== 导出检查 ====================

	printPart("row export check") // 打印测试部分信息

	page.ClickS(page.FindByXPath(actionDropDown)) // 点击操作下拉菜单
	page.ClickS(page.FindByXPath(exportBtn))      // 点击导出按钮
	page.ClickS(page.FindByClass(`grid-batch-1`)) // 点击批量操作按钮

	// ==================== 排序检查 ====================

	printPart("order check") // 打印测试部分信息

	page.Click(idOrderBtn) // 点击 ID 排序按钮
	page.Click(idOrderBtn) // 再次点击 ID 排序按钮

	// ==================== 操作按钮检查 ====================

	printPart("action buttons check") // 打印测试部分信息

	page.Click(rowActionDropDown)   // 点击行操作下拉菜单
	page.Click(previewAction)       // 点击预览操作
	page.Contain("preview content") // 验证页面包含 "preview content" 文本
	page.Display(previewPopup)      // 验证预览弹窗显示

	page.Click(closePreviewAction) // 点击关闭预览操作

	page.Nondisplay(previewPopup) // 验证预览弹窗已隐藏

	page.Click(rowActionDropDown) // 点击行操作下拉菜单
	page.Click(rowAjaxAction)     // 点击行 Ajax 操作

	page.Display(rowAjaxPopup) // 验证行 Ajax 弹窗显示

	page.Click(closeRowAjaxPopup) // 点击关闭行 Ajax 弹窗

	page.Nondisplay(rowAjaxPopup) // 验证行 Ajax 弹窗已隐藏

	wait(2) // 等待 2 秒
}

// ==================== 新建页面操作测试函数 ====================
// testNewPageOperations 测试新建页面的各种操作
// 参数：page - 页面对象，用于操作浏览器
// 说明：该函数测试新建页面的表单填写、选项选择、错误处理、重置、继续创建等功能
func testNewPageOperations(page *Page) {
	page.Click(newPageBtn, 2)                      // 点击新建页面按钮，等待 2 秒
	page.Value(homePageField, "http://google.com") // 验证主页字段值为 "http://google.com"

	// ==================== 选择表单项检查 ====================

	printPart("selections form items check") // 打印测试部分信息

	checkSelectionsInForm(page) // 检查表单中的选择项

	// ==================== 创建错误检查 ====================

	printPart("create error check") // 打印测试部分信息

	page.Click(saveBtn)   // 点击保存按钮
	page.Contain("error") // 验证页面包含 "error" 文本

	// ==================== 重置检查 ====================

	printPart("reset error check") // 打印测试部分信息

	fillNewForm(page, "jane", "girl") // 填写新建表单
	page.Click(resetBtn)              // 点击重置按钮

	// ==================== 继续创建检查 ====================

	printPart("continue creating check") // 打印测试部分信息

	page.Click(inputTab)              // 点击输入标签页
	page.Text(ipField, "")            // 验证 IP 字段为空
	page.Click(continueEditCheckBox)  // 点击继续编辑复选框
	fillNewForm(page, "jane", "girl") // 填写新建表单
	page.Click(saveBtn)               // 点击保存按钮

	// ==================== 创建检查 ====================

	printPart("creating check") // 打印测试部分信息

	fillNewForm(page, "harry", "boy") // 填写新建表单
	page.Click(saveBtn, 2)            // 点击保存按钮，等待 2 秒

	page.NoContain("harry")           // 验证页面不包含 "harry"
	page.Click(genderActionDropDown)  // 点击性别筛选下拉菜单
	page.Click(menOptionActionBtn, 2) // 点击男性选项按钮，等待 2 秒
	page.Click(idOrderBtn)            // 点击 ID 排序按钮
	page.Contain("harry")             // 验证页面包含 "harry"
}

// ==================== 填写新建表单函数 ====================
// fillNewForm 填写新建用户表单
// 参数：
//
//	page - 页面对象，用于操作浏览器
//	name - 用户名
//	gender - 性别
//
// 说明：该函数填写新建用户表单的各个字段
func fillNewForm(page *Page, name, gender string) {
	page.Fill(nameField, name)           // 填写姓名
	page.Fill(ageField, "15")            // 填写年龄
	page.Fill(passwordField, "12345678") // 填写密码
	page.Fill(ipField, "127.0.0.1")      // 填写 IP 地址
	page.Fill(amountField, "15")         // 填写金额
	page.Click(selectTab)                // 点击选择标签页
	page.Click(appleOptField)            // 点击苹果选项
	if gender == "girl" {                // 如果性别为女性
		page.Click(genderGirlCheckBox) // 点击女性复选框
	} else { // 否则
		page.Click(genderBoyCheckBox) // 点击男性复选框
	}
	page.Click(experienceDropDown) // 点击经验下拉菜单
	page.Click(twoYearsSelection)  // 点击两年经验选项
}

// ==================== 检查表单中的选择项函数 ====================
// checkSelectionsInForm 检查表单中的选择项
// 参数：page - 页面对象，用于操作浏览器
// 说明：该函数验证表单中各个选择项的文本和属性是否正确
func checkSelectionsInForm(page *Page) {
	page.Click(selectTab)                       // 点击选择标签页
	page.Text(appleOptField, "Apple")           // 验证苹果选项文本为 "Apple"
	page.Text(bananaOptField, "Banana")         // 验证香蕉选项文本为 "Banana"
	page.Text(watermelonOptField, "Watermelon") // 验证西瓜选项文本为 "Watermelon"
	//page.Text(pearOptField, "") // 验证梨选项（已注释）
	page.Click(experienceDropDown)                                 // 点击经验下拉菜单
	page.Text(twoYearsSelection, "two years")                      // 验证两年经验选项文本
	page.Text(threeYearsSelection, "three years")                  // 验证三年经验选项文本
	page.Text(fourYearsSelection, "four years")                    // 验证四年经验选项文本
	page.Text(fiveYearsSelection, "five years")                    // 验证五年经验选项文本
	page.Click(selectTab)                                          // 点击选择标签页
	page.Attr(page.FindByXPath(multiSelectedOpt), "title", "Beer") // 验证多选已选项的 title 属性为 "Beer"
	page.Click(multiSelectionInput)                                // 点击多选输入框
	page.Text(multiBeerOpt, "Beer")                                // 验证啤酒选项文本为 "Beer"
	page.Text(multiJuiceOpt, "Juice")                              // 验证果汁选项文本为 "Juice"
	page.Text(multiWaterOpt, "Water")                              // 验证水选项文本为 "Water"
	page.Text(multiRedBullOpt, "Red bull")                         // 验证红牛选项文本为 "Red bull"
	page.Click(inputTab)                                           // 点击输入标签页
}

// ==================== 编辑页面操作测试函数 ====================
// testEditPageOperations 测试编辑页面的各种操作
// 参数：page - 页面对象，用于操作浏览器
// 说明：该函数测试编辑页面的表单值显示、编辑、保存等功能
func testEditPageOperations(page *Page) {
	page.Click(rowActionDropDown) // 点击行操作下拉菜单
	page.Click(editPageBtn, 2)    // 点击编辑页面按钮，等待 2 秒

	// ==================== 表单字段值检查 ====================

	printPart("edit form values check") // 打印测试部分信息

	page.Value(nameField, "harry")                   // 验证姓名字段值为 "harry"
	page.Value(homePageField, "http://google.com")   // 验证主页字段值为 "http://google.com"
	page.Value(ageField, "15")                       // 验证年龄字段值为 "15"
	page.Value(emailField, "xxxx@xxx.com")           // 验证邮箱字段值为 "xxxx@xxx.com"
	page.Value(birthdayField, "2010-09-05 00:00:00") // 验证生日字段值为 "2010-09-05 00:00:00"
}
