package html

const ErrorTemplate = `<!DOCTYPE html>
<html>
  <head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>An error occured!</title>
	<link rel="icon" href="/favicon.ico" />
	<style>
	body {
		font-family: Helvetica;
	}
	h1 {
		margin-bottom: 12px;
	}
	</style>
  </head>
  <body>
	<h1>An error occured</h1>
	<code>{{ .Error }}</code>
  </body>
</html>
`
