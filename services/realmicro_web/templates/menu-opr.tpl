{{define "menu-opr"}}
<section>
    {{range .Results.Menus}}
    <a class="list-group-item" href="{{.Url}}">
        <div class="row icon-text {{.Icon}}">
            <p class="btn-title ng-binding">{{.Name}}</p>
        </div>
    </a>
    {{end}}
</section>
{{end}}