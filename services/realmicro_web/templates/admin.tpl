{{define "leftbar"}}{{end}}
{{define "title"}}人员管理{{end}}
{{define "style"}}
{{end}}
{{define "heading"}}<h3>人员管理</h3><hr>{{end}}
{{define "content"}}
<div class="btn-group btn-group-opr"><button type="button" class="btn btn-primary btn-sm" data-toggle="modal" data-target="#addAdminModal">新增</button></div>
<table class="table table-bordered table-striped">
    <thead>
    <th>名字</th>
    <th>中文名</th>
    <th>部门</th>
    <th>角色</th>
    <th>电话</th>
    <th>微信</th>
    <th>创建时间</th>
    <th>操作</th>
    <thead>
    <tbody id="adminList">
    {{range .Results.Admins}}
    <tr id="admin-{{.AdminUser.ID}}">
        <td>{{.AdminUser.Name}}</td>
        <td>{{.AdminUser.NameCn}}</td>
        <td>{{.Department.Name}}({{.Department.NameEn}})</td>
        <td>{{.AdminUser.Role}}</td>
        <td>{{.AdminUser.Phone}}</td>
        <td>{{.AdminUser.Wechat}}</td>
        <td>{{.AdminUser.CreatedAt}}</td>
        <td><button type="button" class="btn btn-danger btn-sm" onclick="delAdmin({{.AdminUser.ID}})">删除</button></td>
    </tr>
    {{end}}
    </tbody>
</table>
<div class="modal fade" id="addAdminModal" tabindex="-1" role="dialog" aria-labelledby="addAdminModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                <h4 class="modal-title">新增人员</h4>
            </div>
            <div class="modal-body">
                <form role="form">
                    <div class="form-group">
                        <label>人员名:</label>
                        <input type="text" class="form-control" id="adminName" placeholder="">
                    </div>
                    <div class="form-group">
                        <label>人员中文名:</label>
                        <input type="text" class="form-control" id="adminNameCn" placeholder="">
                    </div>
                    <div class="form-group">
                        <label>密码:</label>
                        <input type="password" class="form-control" id="adminPassword" placeholder="">
                    </div>
                    <div class="form-group">
                        <label>部门:</label>
                        <select class="form-control" type=text name=department id=department>
                            <option disabled selected> -- 请选择部门 -- </option>
                        {{range .Results.Departments}}
                            <option class="list-group-item" value="{{.ID}}">{{.Name}}({{.NameEn}})</option>
                        {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                        <label>角色:</label>
                        <select class="form-control" type=text name=role id=role>
                            <option disabled selected> -- 请选择角色 -- </option>
                        {{range .Results.Roles}}
                            <option class="list-group-item" value="{{.Role}}">{{.RoleDesc}}</option>
                        {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                        <label>手机号:</label>
                        <input type="text" class="form-control" id="adminPhone" placeholder="">
                    </div>
                    <div class="form-group">
                        <label>微信号:</label>
                        <input type="text" class="form-control" id="adminWechat" placeholder="">
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" onclick="return addAdmin();">提交</button>
            </div>
        </div>
    </div>
</div>
{{end}}
{{define "script"}}
<script>
    function delAdmin(id) {
        swal({
            title: "你确定要删除吗?",
            icon: "warning",
            buttons: true,
            dangerMode: true,
        })
        .then((willDelete) => {
            if (willDelete) {
                $.ajax({
                    url: "/api/v1/del.admin",
                    method: "POST",
                    data: JSON.stringify({"id": id}),
                    success: function (data) {
                        if (data.code !== 0) {
                            toastr.error('删除人员失败');
                        } else {
                            swal("删除人员成功!", {
                                icon: "success",
                            });
                            $('[id="admin-' + id + '"]').remove();
                        }
                    },
                    complete: function () {
                    }
                });
            }
        });
        return false;
    }
    function addAdmin() {
        let name = $("#adminName").val();
        let nameCn = $("#adminNameCn").val();
        let password = $("#adminPassword").val();
        let department = $("#department").val();
        let departmentDesc = $("#department option:selected").text();
        let role = $("#role").val();
        let roleDesc = $("#role option:selected").text();
        let phone = $("#adminPhone").val();
        let wechat = $("#adminWechat").val();
        console.log(name, password, department, $("#department option:selected").text(), role, $("#role option:selected").text(), phone, wechat);
        $("#addName").val("");
        $("#adminPassword").val("");
        $("#department").val("");
        $("#role").val("");
        $("#adminPhone").val("");
        $("#adminWechat").val("");
        $('#addAdminModal').modal('hide');
        $.ajax({
            url: "/api/v1/add.admin",
            method: "POST",
            data: JSON.stringify({
                "name": name,
                "nameCn": nameCn,
                "password": password,
                "departmentId": parseInt(department),
                "role": role,
                "phone": phone,
                "wechat": wechat,
            }),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('新增人员失败');
                } else {
                    toastr.info('新增人员成功');
                    $("#adminList").append(`
                        <tr id="admin-${data.data.id}">
                            <td>${name}</td>
                            <td>${nameCn}</td>
                            <td>${departmentDesc}</td>
                            <td>${roleDesc}</td>
                            <td>${phone}</td>
                            <td>${wechat}</td>
                            <td>${data.data.CreatedAt}</td>
                            <td><button type="button" class="btn btn-danger btn-sm" onclick="delDepartment(${data.data.id})">删除</button></td>
                        </tr>
                    `);
                }
            },
            complete: function () {
            }
        });
        return false;
    }
</script>
{{end}}