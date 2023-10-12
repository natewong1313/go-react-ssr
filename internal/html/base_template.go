package html

const BASE_TEMPLATE = `<!DOCTYPE html>
<html>
  <head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>{{ .Title }}</title>
	{{range $k, $v := .MetaTags}} <meta name="{{$k}}" content="{{$v}}" /> {{end}}
	{{range $k, $v := .OGMetaTags}} <meta property="{{$k}}" content="{{$v}}" /> {{end}}
	{{range .Links}}<link href="{{.Href}}" rel="{{.Rel}}" media="{{.Media}}" hreflang="{{.Hreflang}}" type="{{.Type}}" title="{{.Title}}" />{{end}}
	<link rel="icon" href="/favicon.ico" />
	<style>
	  {{ .CSS }}}
	</style>
  </head>
  <body>
	<div id="root">{{ .ServerHTML }}</div>
	<script>
	  function showError(error) {
		document.getElementById(
		  "root"
		).innerHTML = '<div style="font-family: Helvetica; padding: 4px 16px"> <h1>An error occured</h1> <p style="color: red">'+error+'</p> </div>';
	  }
	</script>
	<script type="module">
	  try{
		{{ .JS }}
	  } catch (e) {
		showError(e.stack)
	  }
	</script>
	{{if .IsDev}}
	<script>
	  let socket = new WebSocket("ws://127.0.0.1:3001/ws");
	  socket.onopen = () => {
		socket.send({{ .RouteID }});
	  };

	  socket.onmessage = (event) => {
		if (event.data === "reload") {
		  console.log("Change detected, reloading...");
		  window.location.reload();
		}
	  };
	</script>
	{{end}}
  </body>
</html>
`
