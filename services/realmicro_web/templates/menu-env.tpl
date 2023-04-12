{{define "menu-env"}}
<section class="panel">
    <header class="panel-heading">
        环境列表
        <span class="pull-right" data-toggle="tooltip" data-placement="bottom" title="通过切换环境、集群来管理不同环境、集群的配置">
                    <img src="/static/img/question.png" class="i-20">
                </span>
    </header>
{{$uri := .Results.Uri}}
{{range .Results.Envs}}
    <div id="treeview" class="no-radius treeview">
        <ul class="list-group">
            <li class="list-group-item node-treeview">
                <span class="icon expand-icon"></span><span class="icon node-icon"></span>
            {{.Env}}
            </li>
        {{range .Envs}}
        {{if .IfSelected }}
        <li class="list-group-item node-treeview node-selected">
        {{else}}
        <li class="list-group-item node-treeview" onclick="self.location='/{{$uri}}?env={{.Env}}&envname={{.Name}}';">
        {{end}}
            <span class="indent"></span>
            <span class="icon glyphicon"></span>
            <span class="icon node-icon"></span>
        {{.Name}}
            <span class="badge">{{.Cluster}}</span>
        </li>
        {{end}}
        </ul>
    </div>
{{end}}
</section>
{{end}}