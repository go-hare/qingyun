{{define "leftbar"}}{{end}}
{{define "title"}}部门管理{{end}}
{{define "heading"}}<h3>部门管理</h3><hr>{{end}}
{{define "content"}}
<div class="btn-group btn-group-opr"><button type="button" class="btn btn-primary btn-sm" data-toggle="modal" data-target="#addDepartmentModal">新增</button></div>
<table class="table table-bordered table-striped">
    <thead>
    <th>部门名</th>
    <th>部门英文名</th>
    <th>创建时间</th>
    <th>操作</th>
    <thead>
    <tbody id="departmentList">
    {{range .Results}}
    <tr id="department-{{.ID}}">
        <td>{{.Name}}</td>
        <td>{{.NameEn}}</td>
        <td>{{.CreatedAt}}</td>
        <td><button type="button" class="btn btn-danger btn-sm" onclick="delDepartment({{.ID}})">删除</button></td>
    </tr>
    {{end}}
    </tbody>
</table>
<div class="modal fade" id="addDepartmentModal" tabindex="-1" role="dialog" aria-labelledby="addDepartmentModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                <h4 class="modal-title">新增部门</h4>
            </div>
            <div class="modal-body">
                <form role="form">
                    <div class="form-group">
                        <label>部门名字:</label>
                        <input type="text" class="form-control" id="addName" placeholder="">
                    </div>
                    <div class="form-group">
                        <label>部门英文名字:</label>
                        <input type="text" class="form-control" id="addNameEn" placeholder="">
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" onclick="return addDepartment();">提交</button>
            </div>
        </div>
    </div>
</div>
{{end}}
{{define "script"}}
    <script>
        function delDepartment(id) {
            swal({
                title: "你确定要删除吗?",
                icon: "warning",
                buttons: true,
                dangerMode: true,
            })
            .then((willDelete) => {
                if (willDelete) {
                    $.ajax({
                        url: "/api/v1/del.department",
                        method: "POST",
                        data: JSON.stringify({"id": id}),
                        success: function (data) {
                            if (data.code !== 0) {
                                toastr.error('删除部门失败');
                            } else {
                                swal("删除部门成功!", {
                                    icon: "success",
                                });
                                $('[id="department-' + id + '"]').remove();
                            }
                        },
                        complete: function () {
                        }
                    });
                }
            });
            return false;
        }
        function addDepartment() {
            let name = $("#addName").val();
            let nameEn = $("#addNameEn").val();
            $("#addName").val("");
            $("#addNameEn").val("");
            $('#addDepartmentModal').modal('hide');
            $.ajax({
                url: "/api/v1/add.department",
                method: "POST",
                data: JSON.stringify({"name": name, "nameEn": nameEn}),
                success: function (data) {
                    if (data.code !== 0) {
                        toastr.error('新增部门失败');
                    } else {
                        toastr.info('新增部门成功');
                        $("#departmentList").append(`
                            <tr id="department-${data.data.id}">
                                <td>${name}</td>
                                <td>${nameEn}</td>
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