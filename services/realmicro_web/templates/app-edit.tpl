{{define "leftbar"}}{{end}}
{{define "title"}}Call{{end}}
{{define "style"}}
pre {
    word-wrap: break-word;
}
{{end}}
{{define "heading"}}<h3>项目管理({{.Results.Service}})</h3><hr>{{end}}
{{define "content"}}
<div class="row">
    <div class="panel panel-default">
        <div class="panel-body">
            <div class="col-sm-12">
                <form id="call-form" onsubmit="return call();">
                    <div class="form-group">
                        <label for="service">部门</label>
                        <ul class="list-group">
                            <select class="form-control" type=text name=department id=department>
                                <option disabled selected> -- 请选择部门 -- </option>
                            {{range .Results.Departments}}
                                <option class="list-group-item" value="{{.ID}}">{{.Name}}({{.NameEn}})</option>
                            {{end}}
                            </select>
                        </ul>
                    </div>
                    <div class="form-group">
                        <label for="request">应用描述</label>
                        <ul class="list-group">
                            <input class="form-control" name=appdesc id=appdesc>
                        </ul>
                    </div>
                    <div class="form-group">
                        <label for="method">应用负责人</label>
                        <ul class="list-group">
                            <select class="form-control" type=text name=owner id=owner>
                                <option disabled selected> -- 请选择应用负责人 -- </option>
                            {{range .Results.Admins}}
                                <option class="list-group-item" value="{{.ID}}">{{.Name}}({{.NameCn}})</option>
                            {{end}}
                            </select>
                        </ul>
                    </div>
                    <div class="form-group">
                        <button class="btn btn-default">提交</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{end}}
{{define "script"}}
<script>
    function call() {
        if (+$("#department").val() === 0) {
            toastr.warning('请选择应用所属部门');
            return false;
        }
        if (+$("#owner").val() === 0) {
            toastr.warning('请选择应用负责人');
            return false;
        }
        var request = {
            "appName": {{.Results.Service}},
            "desc": $("#appdesc").val(),
            "departmentId": +$("#department").val(),
            "owner": +$("#owner").val()
        };
        console.log(request);
        $.ajax({
            url: "/api/v1/edit.app",
            method: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('项目管理提交失败');
                } else {
                    window.location.href="/registry?service={{.Results.Service}}";
                }
            },
            complete: function () {
            }
        });
        return false;
    }
</script>
{{end}}