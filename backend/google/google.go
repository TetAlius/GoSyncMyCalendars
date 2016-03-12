package google

var Config struct {
	googleConfig `json:"google"`
}

type googleConfig struct {
	id          string
	secret      string
	redirectURI string
	loginURI    string
	version     string
	scope       string
}

var Requests struct {
}
