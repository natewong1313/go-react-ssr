package react_renderer

import (
	"bytes"
	"html/template"
)

const HTML_TEMPLATE = `<!DOCTYPE html>
<html>
  <head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>{{ .Title }}</title>
	{{range $k, $v := .MetaTags}}
	<meta name="{{$k}}" content="{{$v}}" />
	{{end}} {{range $k, $v := .OGMetaTags}}
	<meta property="{{$k}}" content="{{$v}}" />
	{{end}}
	<link rel="icon" type="image/svg+xml" href="/react.svg" />
  </head>
  <body>
	<div id="root"></div>
	<style>
	  {{ .CSS }}}
	</style>
	<script>
	  function showError(error) {
		document.getElementById(
		  "root"
		).innerHTML = '<div style="font-family: Helvetica; padding: 4px 16px"> <h1>An error occured</h1> <p style="color: red">${error}</p> </div>';
	  }
	</script>
	<script type="module">
	  try{
		{{ .JS }}
	  } catch (e) {
		showError(e.stack)
	  }
	</script>
	<script>
	  let socket = new WebSocket("ws://127.0.0.1:3001/ws");
	  socket.onopen = () => {
		socket.send({{ .Route }});
	  };

	  socket.onmessage = (event) => {
		if (event.data === "reload") {
		  console.log("Change detected, reloading...");
		  window.location.reload();
		}
	  };

	  // socket.onclose = (event) => {
	  //   socket.send("Client Closed!");
	  // };

	  // socket.onerror = (error) => {};
	</script>
  </body>
</html>
`

type HTMLParams struct {
	Title      string
	MetaTags   map[string]string
	OGMetaTags map[string]string
	JS         template.JS
	CSS        template.CSS
	Route      string
}

func renderHTMLString(params HTMLParams) []byte {
	t := template.Must(template.New("").Parse(HTML_TEMPLATE))
	var output bytes.Buffer
	t.Execute(&output, params)
	return output.Bytes()
}
