package react

type RenderTask struct {
	RouteID           string
	FilePath          string
	Props             string
	RenderConfig      Config
	ServerBuildResult chan ServerBuildResult
	ClientBuildResult chan ClientBuildResult
}

// Build stores the compiled JS and CSS for a React file
type Build struct {
	CompiledJS  string
	CompiledCSS string
}

// Config stores the configuration for a React file to be rendered
type Config struct {
	File     string
	Title    string
	MetaTags map[string]string
	Links    []struct {
		Href     string
		Rel      string
		Media    string
		Hreflang string
		Type     string
		Title    string
	}
	Props interface{}
}

type ClientBuild struct {
	JS           string
	Dependencies []string
}

type ClientBuildResult struct {
	JS           string
	Dependencies []string
	Error        error
}

type ServerRendererBuild struct {
	JS  string
	CSS string
}

type ServerBuildResult struct {
	HTML  string
	CSS   string
	Error error
}
