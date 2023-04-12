{{define "leftbar"}}
<div class="side-bar position-fixed">
    <div class="J_appFound">
        {{template "menu-env" .}}

        {{template "menu-opr" .}}
    </div>
</div>
{{end}}
{{define "heading"}}<h4><input class="form-control input-lg search" type=text placeholder="Search"/></h4>{{end}}
{{define "title"}}Registry{{end}}
{{define "content"}}
<div>
{{$env := .Results.Env}}
{{$envName := .Results.EnvName}}
{{range .Results.Projects}}
    <h3 style="margin-top: 20px;"><span id="project-{{.Project}}">{{.ProjectAlias}}</span> <button type="button" class="btn btn-default" onclick="return showProjectModal('{{.Project}}');">
        <span class="glyphicon glyphicon-edit" aria-hidden="true"></span>
    </button></h3>
    <hr>
    {{range .Services}}
        <a href="registry?service={{.Name}}&env={{$env}}&envname={{$envName}}" data-filter="{{.Name}}" class="btn btn-default btn-lg" style="margin: 5px 3px 5px 3px;">{{.Name}}</a>
    {{end}}
{{end}}
</div>
<div class="modal fade" id="modifyProjectModal" tabindex="-1" role="dialog" aria-labelledby="modifyProjectModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                <h4 class="modal-title">修改项目名 <span id="projectLabel"></span></h4>
            </div>
            <div class="modal-body">
                <form role="form">
                    <div class="form-group">
                        <input type="text" class="form-control" id="projectAlias" placeholder="请输入项目名称">
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" onclick="return projectAlias();">提交更改</button>
            </div>
        </div>
    </div>
</div>
{{end}}
{{define "script"}}
<script type="text/javascript">
    jQuery(function($, undefined) {
        var refs = $('a[data-filter]');
        $('.search').on('keyup', function() {
            var val = $.trim(this.value);
            refs.hide();
            refs.filter(function() {
                return $(this).data('filter').search(val) >= 0
            }).show();
        });
    });

    function showProjectModal(project) {
        $("#projectLabel").html(project);
        $("#modifyProjectModal").modal();
    }

    function projectAlias() {
        let project = $("#projectLabel").html();
        let alias = $("#projectAlias").val();
        $("#projectAlias").val("");
        $('#modifyProjectModal').modal('hide');
        $.ajax({
            url: "/api/v1/edit.project",
            method: "POST",
            data: JSON.stringify({"project": project, "alias": alias}),
            success: function (data) {
                if (data.code !== 0) {
                    toastr.error('设置项目名字失败');
                } else {
                    toastr.info('设置项目名字成功');
                    $("#project-"+project).html(alias+"("+project+")");
                }
            },
            complete: function () {
            }
        });
        return false;
    }
</script>
{{end}}