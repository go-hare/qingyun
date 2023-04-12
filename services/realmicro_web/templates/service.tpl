{{define "head"}}
<link rel="stylesheet" href="/static/css/jquery.jsonview.css">
<link rel="stylesheet" href="/static/css/tree.css">
{{end}}
{{define "leftbar"}}
<div class="side-bar position-fixed">
    <div class="J_appFound">

        <section class="panel">
            <header class="panel-heading">
                项目信息
                <span class="pull-right" data-toggle="tooltip" data-placement="bottom" title="" data-original-title="修改项目基本信息">
                    <a href="/app-edit?service={{.Results.AppDetail.App.AppName}}">
                        <img src="/static/img/edit.png" class="i-20">
                    </a>
                </span>
            </header>
            <div class="panel-body">
                <table class="project-info">
                    <tbody class="text-left">
                    <tr>
                        <th>应用名:</th>
                        <td>
                            {{.Results.AppDetail.App.AppName}}
                        </td>
                    </tr>
                    <tr>
                        <th>部门:</th>
                        <td>{{if .Results.AppDetail.Department.Name}} {{.Results.AppDetail.Department.Name}} {{else}} - {{end}}</td>
                    </tr>
                    <tr>
                        <th>负责人:</th>
                        <td>{{if .Results.AppDetail.AdminUser.Name}} {{.Results.AppDetail.AdminUser.Name}} {{else}} - {{end}}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </section>

        <!--operation entrance-->
        <section>
            <a class="list-group-item hover ng-isolate-scope" href="/app-edit?service={{.Results.AppDetail.App.AppName}}">
                <div class="row icon-text icon-project-manage">
                    <p class="btn-title ng-binding">管理项目</p>
                </div>
            </a>

        <!-- <a class="list-group-item ng-hide">
                <div class="row icon-text icon-plus-orange">
                    <p class="btn-title ng-binding">补缺环境</p>
                </div>
            </a>-->
        </section>

    </div>
</div>
{{end}}

{{define "title"}}Service{{end}}
{{define "heading"}}<h3>{{with $svc := index .Results.Services 0}}{{$svc.Name}}{{end}}</h3>{{end}}
{{define "content"}}
    <hr>
    <ul id="myTab" class="nav nav-tabs">
        <li class="active"><a href="#home" data-toggle="tab" class="nav-tab-normal">Nodes</a></li>
        <li><a href="#endpoints" data-toggle="tab" class="nav-tab-normal">Endpoints</a></li>
        <li><a href="#config" data-toggle="tab" class="nav-tab-normal">Config</a></li>
    </ul>
    <div id="myTabContent" class="tab-content">
        <div class="tab-pane fade in active" id="home">
            <hr>
            {{range .Results.Services}}
                <h5>Version {{.Version}}, Nodes Num {{len .Nodes}}</h5>
                <table class="table table-bordered table-striped">
                    <thead>
                    <th>Id</th>
                    <th>Address</th>
                    <th>Port</th>
                    <th>Metadata</th>
                    <thead>
                    <tbody>
                    {{range .Nodes}}
                    <tr>
                        <td>{{.Id}}</td>
                        <td>{{.Address}}</td>
                        <td>{{.Port}}</td>
                        <td>{{ range $key, $value := .Metadata }}{{$key}}={{$value}} {{end}}</td>
                    </tr>
                    {{end}}
                    </tbody>
                </table>
            {{end}}
        </div>
        <div class="tab-pane fade" id="endpoints">
            {{with $svc := index .Results.Services 0}}
            <hr/>
            {{range $i, $v := $svc.Endpoints}}
                <div class="panel panel-default">
                    <div class="panel-heading">
                        <h4 class="panel-title">
                            <a data-toggle="collapse" data-parent="#accordion"
                               href="#collapse{{$i}}">
                                <h4>{{$v.Name}}</h4>
                            </a>
                        </h4>
                    </div>
                    <div id="collapse{{$i}}" class="panel-collapse collapse">
                        <div class="panel-body">
                            <table class="table table-bordered">
                                <tbody>
                                <tr>
                                    <th class="col-sm-2" scope="row">Metadata</th>
                                    <td>{{ range $key, $value := $v.Metadata }}{{$key}}={{$value}} {{end}}</td>
                                </tr>
                                <tr>
                                    <th class="col-sm-2" scope="row">Request</th>
                                    <td><pre>{{format $v.Request}}</pre></td>
                                </tr>
                                <tr>
                                    <th class="col-sm-2" scope="row">Response</th>
                                    <td><pre>{{format $v.Response}}</pre></td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            {{end}}
            {{end}}
        </div>
        <div class="tab-pane fade" id="config">
            <hr>
            <div class="row">
                <div class="col-md-6">
                    <div class="tree prop-block">
                        <ul id="">
                            <li id="liNode-">
                                <div class="node" id="dirNode-" onclick="configValues('')">{{.Results.Svc}}</div>
                                <div class="btn-group ng-scope">
                                    <button class="btn btn-default btn-xs btn-list" type="button" title="Create Directory" data-toggle="modal" data-target="#addConfigDirModal">
                                        <span class="glyphicon glyphicon-plus "></span>
                                    </button>
                                </div>
                                {{fc .Results.Config}}
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="prop-block ng-scope">
                        <button type="button" class="btn btn-default btn-xs ng-scope" title="增加配置项" data-toggle="modal" data-target="#addConfigItemModal">
                            增加配置项
                        </button>
                        <table class="property-table ng-scope">
                            <tbody id="config-list">
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade" id="addConfigDirModal" tabindex="-1" role="dialog" aria-labelledby="addConfigDirModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    <h4 class="modal-title">增加配置目录</h4>
                </div>
                <div class="modal-body">
                    <form role="form">
                        <div class="form-group">
                            <input type="text" class="form-control" id="addConfigDirValue" placeholder="">
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                    <button type="button" class="btn btn-primary" onclick="return addConfigDir();">提交</button>
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade" id="editConfigModal" tabindex="-1" role="dialog" aria-labelledby="editConfigModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    <h4 class="modal-title">配置项 <span id="configValueLabel"></span></h4>
                </div>
                <div class="modal-body">
                    <form role="form">
                        <div class="form-group">
                            <textarea class="form-control" style="height: 300px;" id="configValueEdit" placeholder=""></textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                    <button type="button" class="btn btn-primary" onclick="return editConfig();">提交</button>
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade" id="addConfigItemModal" tabindex="-1" role="dialog" aria-labelledby="addConfigItemModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    <h4 class="modal-title">增加配置项 <span id="addConfigItemLabel"></span></h4>
                </div>
                <div class="modal-body">
                    <form role="form">
                        <div class="form-group">
                            <label>Key:</label>
                            <input type="text" class="form-control" id="addConfigItemKey" placeholder="">
                        </div>
                        <div class="form-group">
                            <label>Value:</label>
                            <textarea class="form-control" style="height: 300px;" id="addConfigItemValue" placeholder=""></textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                    <button type="button" class="btn btn-primary" onclick="return addConfigItem();">提交</button>
                </div>
            </div>
        </div>
    </div>
{{end}}
{{define "script"}}
<script src="/static/js/jquery.jsonview.js"></script>
<script src="/static/js/jquery.json-editor.min.js"></script>
<script>
    let svc = "{{.Results.Svc}}";
    let cj = {{.Results.Config}};
    let nowSelectedDir = "";
    let env = {{.Results.Env}};
    let envName = {{.Results.EnvName}};

    function isJsonString(str) {
        try {
            if (typeof JSON.parse(str) === "object") {
                return true;
            }
        } catch(e) {
        }
        return false;
    }

    function configValues(key) {
        nowSelectedDir = key;
        console.log("now selected dir:", nowSelectedDir);
        $(".active-node").removeClass("active-node");
        $(".btn-group-show").removeClass("btn-group-show");
        $('[id="dirNode-' + key + '"]').addClass("active-node");
        $('[id="dirNodeBtn-' + key + '"]').addClass("btn-group-show");
        let keys = key.split("/");
        let data = cj;
        if (nowSelectedDir !== "") {
            for (let i = 0; i < keys.length; i++) {
                data = data[keys[i]].nodes;
            }
        }
        $("#config-list").empty();
        for (let k in data) {
            if (data[k].dir) {
                continue;
            }
            cvKey = data[k].key.replace(new RegExp("\\.","g"),"-");
            $("#config-list").append(`
                <tr class="ng-scope odd-row" id="config-${svc}/${data[k].longKey}">
                    <td>
                        <div class="btn-group ng-scope btn-group-td">
                            <button type="button" class="btn btn-default btn-xs ng-scope btn-list" title="删除配置" onclick="delConfig('${svc}/${data[k].longKey}');">
                                <span class="glyphicon glyphicon-trash"></span>
                            </button>
                            <button type="button" class="btn btn-default btn-xs ng-scope btn-list" onclick="editConfigModel('${data[k].longKey}');" title="编辑">
                                <span class="glyphicon glyphicon-edit"></span>
                            </button>
                        </div>
                    </td>
                    <td style="text-align: left; font-weight: bold; padding: 5px 5px;">
                        <span class="ng-binding">${data[k].key}:</span>
                    </td>
                    <td style="padding: 5px 5px;">
                        <pre class="value ng-scope" style="margin: 0 0;" id="configValue-${cvKey}"></pre>
                    </td>
                </tr>
            `);
            console.log("data value:", data[k].value);
            if (isJsonString(data[k].value)) {
                console.log("config json id:", "#configValue-"+cvKey);
                $("#configValue-"+cvKey).JSONView(JSON.parse(data[k].value), { collapsed: false });
            } else {
                console.log(data[k].value, ": is not json");
                $("#configValue-"+cvKey).html(data[k].value);
            }
        }
        return false;
    }
    function delConfig(key) {
        console.log("del key:", key);
        swal({
            title: "你确定要删除吗?",
            text: "删除后配置将被清空!",
            icon: "warning",
            buttons: true,
            dangerMode: true,
        })
        .then((willDelete) => {
            if (willDelete) {
                $.ajax({
                    url: "/api/v1/del.config",
                    method: "POST",
                    data: JSON.stringify({"env": env, "envName": envName, "key": key}),
                    success: function (data) {
                        if (data.code !== 0) {
                            toastr.error('删除配置失败');
                        } else {
                            swal("删除配置成功!", {
                                icon: "success",
                            });
                            $('[id="config-' + key + '"]').remove();
                        }
                    },
                    complete: function () {
                    }
                });
            }
        });
        return false;
    }
    function delConfigDir(key) {
        swal({
            title: "你确定要删除吗?",
            text: "删除目录后，目录下的所有配置将被清空!",
            icon: "warning",
            buttons: true,
            dangerMode: true,
        })
        .then((willDelete) => {
            if (willDelete) {
                $.ajax({
                    url: "/api/v1/del.config",
                    method: "POST",
                    data: JSON.stringify({"env": env, "envName": envName, "key": svc+"/"+key}),
                    success: function (data) {
                        if (data.code !== 0) {
                            toastr.error('删除目录配置失败');
                        } else {
                            swal("删除目录配置成功!", {
                                icon: "success",
                            });
                            $('[id="ulNode-' + key + '"]').remove();
                        }
                    },
                    complete: function () {
                    }
                });
            }
        });
    }
    function editConfigModel(key) {
        console.log("edit key:", key);
        let keys = key.split("/");
        let data = cj;
        for (let i = 0; i < keys.length; i++) {
            if (data[keys[i]].nodes) {
                data = data[keys[i]].nodes;
            } else {
                data = data[keys[i]]
            }
        }
        console.log("data:", data);
        $("#configValueLabel").html(key);
        $("#configValueEdit").val(data.value);
        $("#editConfigModal").modal();
    }
    function editConfig() {
        let key = $("#configValueLabel").html();
        let value = $("#configValueEdit").val();
        console.log("key:", key);
        console.log("value:", value);
        $('#editConfigModal').modal('hide');
        $.ajax({
            url: "/api/v1/put.config",
            method: "POST",
            data: JSON.stringify({"env": env, "envName": envName, "key": svc+"/"+key, "value": value}),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('编辑配置失败');
                } else {
                    toastr.info('编辑配置成功');
                    let keys = key.split("/");
                    let data = cj;
                    for (let i = 0; i < keys.length; i++) {
                        if (data[keys[i]].nodes) {
                            data = data[keys[i]].nodes;
                        } else {
                            data = data[keys[i]]
                        }
                    }
                    data.value = value;
                    if (isJsonString(value)) {
                        $("#configValue-"+keys[keys.length-1]).JSONView(JSON.parse(value), { collapsed: false });
                    } else {
                        $("#configValue-"+keys[keys.length-1]).html(value);
                    }
                }
            },
            complete: function () {
            }
        });
        return false;
    }
    function addConfigDir() {
        let key = $("#addConfigDirValue").val();
        console.log("add config dir:", key);
        $("#addConfigDirValue").val("");
        $('#addConfigDirModal').modal('hide');
        let fullKey = svc+"/"+nowSelectedDir+"/"+key;
        let pathKey = nowSelectedDir+"/"+key;
        if (nowSelectedDir === "") {
            fullKey = svc+"/"+key;
            pathKey = key
        }
        $.ajax({
            url: "/api/v1/put.config",
            method: "POST",
            data: JSON.stringify({"env": env, "envName": envName, "key": fullKey, "ifDir": true}),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('新增配置目录失败');
                } else {
                    toastr.info('新增配置目录成功');
                    let data = cj;
                    if (nowSelectedDir !== "") {
                        let keys = nowSelectedDir.split("/");
                        for (let i = 0; i < keys.length; i++) {
                            if (data[keys[i]].nodes) {
                                data = data[keys[i]].nodes;
                            } else {
                                data = data[keys[i]]
                            }
                        }
                    }
                    data[key] = {"key": key, "longKey": pathKey, "dir": true, value: "etcdv3_dir_$2H#%gRe3*t"};
                    $('[id="liNode-' + nowSelectedDir + '"]').append(`
                        <ul id="ulNode-${data[key].longKey}"><li class="ng-scope" id="liNode-${data[key].longKey}"><div class="node" id="dirNode-${data[key].longKey}" onclick="configValues('${data[key].longKey}')">${key}</div>
                            <div class="btn-group ng-scope btn-group-hide" id="dirNodeBtn-${data[key].longKey}">
                                <button class="btn btn-default btn-xs ng-scope btn-list" type="button" title="Create Directory" data-toggle="modal" data-target="#addConfigDirModal">
                                    <span class="glyphicon glyphicon-plus"></span>
                                </button>
                                <button type="button" class="btn btn-default btn-xs ng-scope btn-list" title="Delete Directory" onclick="delConfigDir('${data[key].longKey}')">
                                    <span class="glyphicon glyphicon-trash"></span>
                                </button>
                            </div>
                        </li></ul>
                    `);
                }
            },
            complete: function () {
            }
        });
        return false;
    }
    function addConfigItem() {
        let key = $("#addConfigItemKey").val();
        let value = $("#addConfigItemValue").val();
        $("#addConfigItemKey").val("");
        $("#addConfigItemValue").val("");
        $('#addConfigItemModal').modal('hide');
        let fullKey = svc+"/"+nowSelectedDir+"/"+key;
        let pathKey = nowSelectedDir+"/"+key;
        if (nowSelectedDir === "") {
            fullKey = svc+"/"+key;
            pathKey = key
        }
        $.ajax({
            url: "/api/v1/put.config",
            method: "POST",
            data: JSON.stringify({"env": env, "envName": envName, "key": fullKey, "value": value}),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('新增配置项失败');
                } else {
                    toastr.info('新增配置项成功');
                    let data = cj;
                    if (nowSelectedDir !== "") {
                        let keys = nowSelectedDir.split("/");
                        for (let i = 0; i < keys.length; i++) {
                            if (data[keys[i]].nodes) {
                                data = data[keys[i]].nodes;
                            } else {
                                data = data[keys[i]]
                            }
                        }
                    }
                    data[key] = {"key": key, "longKey": pathKey, value: value};
                    $("#config-list").append(`
                        <tr class="ng-scope odd-row" id="config-${svc}/${data[key].longKey}">
                            <td>
                                <div class="btn-group ng-scope">
                                    <button type="button" class="btn btn-default btn-xs ng-scope btn-list" title="删除配置" onclick="delConfig('${svc}/${data[key].longKey}');">
                                        <span class="glyphicon glyphicon-trash"></span>
                                    </button>
                                    <button type="button" class="btn btn-default btn-xs ng-scope btn-list" onclick="editConfigModel('${data[key].longKey}');" title="编辑">
                                        <span class="glyphicon glyphicon-edit"></span>
                                    </button>
                                </div>
                            </td>
                            <td style="text-align: left; font-weight: bold; padding: 5px 5px;">
                                <span class="ng-binding">${data[key].key}:</span>
                            </td>
                            <td style="padding: 5px 5px;">
                                <pre class="value ng-scope" style="margin: 0 0;" id="configValue-${data[key].key}"></pre>
                            </td>
                        </tr>
                    `);
                    if (isJsonString(data[key].value)) {
                        $("#configValue-"+data[key].key).JSONView(JSON.parse(data[key].value), { collapsed: false });
                    } else {
                        $("#configValue-"+data[key].key).html(data[key].value);
                    }
                }
            },
            complete: function () {
            }
        });
        return false
    }
</script>
{{end}}