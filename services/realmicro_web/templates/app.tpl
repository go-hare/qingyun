{{define "leftbar"}}{{end}}
{{define "title"}}服务管理{{end}}
{{define "style"}}
{{end}}
{{define "heading"}}<h3>服务管理</h3><hr>{{end}}
{{define "content"}}
<ul id="myTab" class="nav nav-tabs">
    <li class="active"><a href="#myApp" data-toggle="tab" class="nav-tab-normal">我的APP</a></li>
    <li><a href="#allApp" data-toggle="tab" class="nav-tab-normal">所有APP</a></li>
</ul>
<div id="myTabContent" class="tab-content">
    <div class="tab-pane fade in active" id="myApp">
        <table class="table table-bordered table-striped">
            <thead>
            <th>服务名</th>
            <th>描述</th>
            <th>所属部门</th>
            <th>负责人</th>
            <th>操作</th>
            <thead>
            <tbody id="adminList">
            {{range .Results.MyApps}}
            <tr id="">
                <td>{{.App.AppName}}</td>
                <td>{{.App.Desc}}</td>
                <td>{{.Department.Name}}({{.Department.NameEn}})</td>
                <td>{{.AdminUser.Name}}({{.AdminUser.NameCn}})</td>
                <td><button type="button" class="btn btn-primary btn-sm" onclick="self.location='/app-edit?service={{.App.AppName}}'">编辑</button></td>
            </tr>
            {{end}}
            </tbody>
        </table>
    </div>
    <div class="tab-pane fade" id="allApp">
        <table class="table table-bordered table-striped">
            <thead>
            <th>服务名</th>
            <th>描述</th>
            <th>所属部门</th>
            <th>负责人</th>
            <th>操作</th>
            <thead>
            <tbody id="adminList">
            {{range .Results.Apps}}
            <tr id="">
                <td>{{.App.AppName}}</td>
                <td>{{.App.Desc}}</td>
                <td>{{.Department.Name}}({{.Department.NameEn}})</td>
                <td>{{.AdminUser.Name}}({{.AdminUser.NameCn}})</td>
                <td><button type="button" class="btn btn-primary btn-sm" onclick="self.location='/app-edit?service={{.App.AppName}}'">编辑</button></td>
            </tr>
            {{end}}
            </tbody>
        </table>
    </div>
</div>
{{end}}
{{define "script"}}
<script>
    function detailApp() {
        return false;
    }
</script>
{{end}}