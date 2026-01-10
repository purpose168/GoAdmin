package types

// tmpls 是模板映射表，包含各种前端交互模板
// 用于处理表单字段的联动、选择、显示/隐藏等交互逻辑
var tmpls = map[string]string{"choose": `{{define "choose"}}
    <script>
        // 监听选择框的选择事件
        $("select.{{.Field}}").on("select2:select", function (e) {
            // 当选中的值等于指定值时，更新关联字段的值
            if (e.params.data.text === {{.Val}} || e.params.data.id === {{.Val}}) {
                // 如果关联字段是选择框，则设置选择值
                if ($("select.{{.ChooseField}}").length > 0) {
                    $("select.{{.ChooseField}}").val("{{.Value}}[0]").select2()
                } else {
                    // 否则直接设置值
                    $(".{{.ChooseField}}").val({{.Value}})
                }
            }
        })
    </script>
{{end}}`, "choose_ajax": `{{define "choose_ajax"}}
    <script>
        // 更新双列表框的选项
        let {{.Field}}_updateBoxSelections = function (selectObj, new_opts) {
            selectObj.html('');
            new_opts.forEach(function (opt) {
                selectObj.append($('<option value="' + opt["id"] + '">' + opt["text"] + '</option>'));
            });
            selectObj.bootstrapDualListbox('refresh', true);
        };

        // 发送AJAX请求获取关联字段的数据
        let {{.Field}}_req = function (selectObj, box, event) {
            $.ajax({
                url: "{{.Url}}", // 请求URL
                type: 'post', // 请求类型
                dataType: 'text', // 数据类型
                data: {
                    'value': $("select.{{.Field}}").val(), // 当前字段的值
                    {{.PassValue}} // 传递的额外参数
                    'event': event // 事件类型
                },
                success: function (data) {
                    // 解析返回的数据
                    if (typeof (data) === "string") {
                        data = JSON.parse(data);
                    }
                    if (data.code === 0) {
                        {{if eq .ActionJS ""}}
                        // 如果没有自定义动作脚本，执行默认逻辑
                        if (selectObj.length > 0) {
                            if (typeof (data.data) === "object") {
                                // 如果是双列表框
                                if (box) {
                                    {{.Field}}_updateBoxSelections(selectObj, data.data)
                                } else {
                                    // 如果是多选框，清空选项
                                    if (typeof (selectObj.attr("multiple")) !== "undefined") {
                                        selectObj.html("");
                                    }
                                    // 更新选择框选项
                                    selectObj.select2({
                                        data: data.data
                                    });
                                }
                            } else {
                                // 如果是单选框
                                if (box) {
                                    selectObj.val(data.data).select2()
                                } else {

                                }
                            }
                        } else {
                            // 更新普通字段的值
                            $('.{{.ChooseField}}').val(data.data);
                        }

                        {{else}}
                        // 执行自定义动作脚本
                        {{.ActionJS}}

                        {{end}}
                    } else {
                        // 显示错误消息
                        swal(data.msg, '', 'error');
                    }
                },
                error: function () {
                    // 显示错误提示
                    alert('error')
                }
            });
        };

        // 如果不是双列表框
        if ($("label[for='{{.Field}}']").next().find(".bootstrap-duallistbox-container").length === 0) {
            // 监听选择事件
            $("select.{{.Field}}").on("select2:select", function (e) {
                let id = '{{.ChooseField}}';
                let selectObj = $("select." + id);
                // 清空关联字段
                if (selectObj.length > 0) {
                    selectObj.val("").select2();
                    selectObj.html('<option value="" selected="selected"></option>')
                }
                // 发送请求获取关联字段的数据
                {{.Field}}_req(selectObj, false, "select");
            });
            // 如果是多选框，监听取消选择事件
            if (typeof ($("select.{{.Field}}").attr("multiple")) !== "undefined") {
                $("select.{{.Field}}").on("select2:unselect", function (e) {
                    let id = '{{.ChooseField}}';
                    let selectObj = $("select." + id);
                    // 清空关联字段
                    if (selectObj.length > 0) {
                        selectObj.val("").select2();
                        selectObj.html('<option value="" selected="selected"></option>')
                    }
                    // 发送请求获取关联字段的数据
                    {{.Field}}_req(selectObj, false, "unselect");
                })
            }
        } else {
            // 双列表框的处理逻辑
            let {{.Field}}_lastState = $(".{{.Field}}").val();

            // 监听变化事件
            $(".{{.Field}}").on('change', function (e) {
                var newState = $(this).val();
                // 检查是否有选项被取消选择
                if ($({{.Field}}_lastState).not(newState).get().length > 0) {
                    let id = '{{.ChooseField}}';
                    {{.Field}}_req($("." + id), true, "unselect");
                }
                // 检查是否有新选项被选择
                if ($(newState).not({{.Field}}_lastState).get().length > 0) {
                    let id = '{{.ChooseField}}';
                    {{.Field}}_req($("." + id), true, "select");
                }
                {{.Field}}_lastState = newState;
            })
        }
    </script>
{{end}}`, "choose_custom": `{{define "choose_custom"}}
    <script>
        // 监听选择框的选择事件，执行自定义JavaScript代码
        $("select.{{.Field}}").on("select2:select", function (e) {
            {{.JS}}
        })
    </script>
{{end}}`, "choose_disable": `{{define "choose_disable"}}
    <script>
        // 监听选择框的选择事件，根据选中的值禁用或启用关联字段
        $("select.{{.Field}}").on("select2:select", function (e) {
            // 如果选中的值在指定值列表中，禁用关联字段
            if ({{.Value}}.indexOf(e.params.data.text) !== -1 || {{.Value}}.indexOf(e.params.data.id) !== -1) {
                {{range $key, $fields := .ChooseFields}}

                // 禁用关联字段
                $(".{{$fields}}").prop('disabled', true);

                {{end}}
            } else {
                {{range $key, $fields := .ChooseFields}}

                // 启用关联字段
                $(".{{$fields}}").prop('disabled', false);

                {{end}}
            }
        });
    </script>
{{end}}`, "choose_hide": `{{define "choose_hide"}}
    <script>
        // 监听选择框的选择事件，根据选中的值隐藏或显示关联字段
        $("select.{{.Field}}").on("select2:select", function (e) {
            // 如果选中的值在指定值列表中，隐藏关联字段
            if ({{.Value}}.indexOf(e.params.data.text) !== -1 || {{.Value}}.indexOf(e.params.data.id) !== -1) {
                {{range $key, $fields := .ChooseFields}}

                // 隐藏关联字段
                $("label[for='{{$fields}}']").parent().hide();

                {{end}}
            } else {
                {{range $key, $fields := .ChooseFields}}

                // 显示关联字段
                $("label[for='{{$fields}}']").parent().show();

                {{end}}
            }
        });
        // 页面加载时初始化关联字段的显示状态
        $(function () {
            // 获取当前字段的选中值
            let {{.Field}}data = $(".{{.Field}}").select2("data");
            let {{.Field}}text = "";
            let {{.Field}}id = "";
            if ({{.Field}}data.length > 0) {
                {{.Field}}text = {{.Field}}data[0].text;
                {{.Field}}id = {{.Field}}data[0].id;
            }
            // 如果当前选中的值在指定值列表中，隐藏关联字段
            if ({{.Value}}.indexOf({{$.Field}}text) !== -1 || {{.Value}}.indexOf({{$.Field}}id) !== -1) {
                {{range $key, $fields := .ChooseFields}}

                // 隐藏关联字段
                $("label[for='{{$fields}}']").parent().hide();

                {{end}}
            }
        })
    </script>
{{end}}`, "choose_map": `{{define "choose_map"}}
    <script>
        // 监听选择框的选择事件，根据映射关系执行不同的操作
        $("select.{{.Field}}").on("select2:select", function (e) {
            {{range $val, $object := .Data}}

            {{if $object.Hide}}
            // 如果需要隐藏字段
            if (e.params.data.text === "{{$val}}" || e.params.data.id === "{{$val}}") {
                // 隐藏关联字段
                $("label[for='{{$object.Field}}']").parent().hide()
            } else {
                // 显示关联字段
                $("label[for='{{$object.Field}}']").parent().show()
            }

            {{else if $object.Disable}}
            // 如果需要禁用字段
            if (e.params.data.text === "{{$val}}" || e.params.data.id === "{{$val}}") {
                // 禁用关联字段
                $("label[for='{{$object.Field}}']").prop('disabled', true);
            } else {
                // 启用关联字段
                $("label[for='{{$object.Field}}']").prop('disabled', false);
            }

            {{else}}
            // 如果需要设置字段值
            if (e.params.data.text === "{{$val}}" || e.params.data.id === "{{$val}}") {
                // 如果关联字段是选择框，则设置选择值
                if ($("select.{{$object.Field}}").length > 0) {
                    $("select.{{$object.Field}}").val("{{$object.Value}}").select2()
                } else {
                    // 否则直接设置值
                    $("#{{$object.Field}}").val("{{$object.Value}}")
                }
            }

            {{end}}

            {{end}}
        })
    </script>
{{end}}`, "choose_show": `{{define "choose_show"}}
    <script>
        // 监听选择框的选择事件，根据选中的值显示或隐藏关联字段
        $("select.{{.Field}}").on("select2:select", function (e) {
            // 如果选中的值在指定值列表中，显示关联字段
            if ({{.Value}}.indexOf(e.params.data.text) !== -1 || {{.Value}}.indexOf(e.params.data.id) !== -1) {
                {{range $key, $fields := .ChooseFields}}

                // 显示关联字段
                $("label[for='{{$fields}}']").parent().show();

                {{end}}
            } else {
                {{range $key, $fields := .ChooseFields}}

                // 隐藏关联字段
                $("label[for='{{$fields}}']").parent().hide();

                {{end}}
            }
        });
        // 页面加载时初始化关联字段的显示状态
        $(function () {
            // 获取当前字段的选中值
            let {{.Field}}data = $(".{{.Field}}").select2("data");
            let {{.Field}}text = "";
            let {{.Field}}id = "";
            if ({{.Field}}data.length > 0) {
                {{.Field}}text = {{.Field}}data[0].text;
                {{.Field}}id = {{.Field}}data[0].id;
            }
            // 如果当前选中的值在指定值列表中，显示关联字段
            if ({{.Value}}.indexOf({{$.Field}}text) !== -1 || {{.Value}}.indexOf({{$.Field}}id) !== -1) {
                {{range $key, $fields := .ChooseFields}}

                // 显示关联字段
                $("label[for='{{$fields}}']").parent().show();

                {{end}}
            }
        })
    </script>
{{end}}`}
