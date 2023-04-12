{{define "layout"}}
<html>
<head>
    <title>REAL PaaS Cloud</title>
    <link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/3.3.7/css/bootstrap.min.css" crossorigin="anonymous">
    <link rel="stylesheet" href="/static/css/common-style.css">
    <link rel="stylesheet" href="https://lf26-cdn-tos.bytecdntp.com/cdn/expire-1-M/toastr.js/2.1.4/toastr.min.css">
    <style>
        {{ template "style" . }}
    </style>
    {{ template "head" . }}
</head>
<body>
<nav class="navbar navbar-inverse">
    <div class="container">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navBar">
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="/">REAL PaaS Cloud</a>
        </div>
        <div class="collapse navbar-collapse" id="navBar">
            <ul class="nav navbar-nav navbar-right">
                <li class="navbar-right-label">{{.LoginUser}}</li>
                <li class="navbar-right-label"><button type="button" class="btn btn-primary btn-xs navbar-right-btn" onclick="self.location='/logout'">登出</button></li>
                {{/*<li><a href="call">Call</a></li>*/}}
                {{/*{{if .StatsURL}}<li><a href="{{.StatsURL}}" class="navbar-link">Stats</a></li>{{end}}*/}}
            </ul>
        </div>
    </div>
</nav>

{{ template "leftbar" . }}

<div class="container">
    <div class="row">
        <div class="col-sm-12">
            {{ template "heading" . }}
            {{ template "content" . }}
        </div>
    </div>
</div>

<div class="container">
    <div class="footer footer-copyright">
        Copyright &copy; 2016 - <span id="year"></span> RealTech-inc.
    </div>
</div>

<script src="https://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/3.3.7/js/bootstrap.min.js" crossorigin="anonymous"></script>
<script src="https://lf6-cdn-tos.bytecdntp.com/cdn/expire-1-M/toastr.js/2.1.4/toastr.min.js"></script>
<script src="https://lf3-cdn-tos.bytecdntp.com/cdn/expire-1-M/sweetalert/2.1.2/sweetalert.min.js"></script>
<script src="/static/js/is.min.js"></script>
{{template "script" . }}

<script>
    var d = new Date(); document.getElementById('year').innerHTML = d.getFullYear();
    $(function () { $("[data-toggle='tooltip']").tooltip(); });
</script>
</body>
</html>
{{end}}
{{ define "style" }}{{end}}
{{ define "head" }}{{end}}
{{ define "leftbar" }}{{end}}
{{ define "script" }}{{end}}
{{ define "title" }}{{end}}
{{ define "heading" }}<h3>&nbsp;</h3>{{end}}