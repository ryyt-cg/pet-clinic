package info

// response
// Application Info response
type response struct {
	AppName     string `json:"appName"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Ip          string `json:"ip"`
	GitCommit   string `json:"gitCommit"`
}
