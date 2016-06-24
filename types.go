package main

// Params represents the valid paramenter options for the webhook plugin.
type Params struct {
	URLs        []string          `json:"urls"`
	Debug       bool              `json:"debug"`
	Auth        Auth              `json:"auth"`
	Headers     map[string]string `json:"header"`
	Method      string            `json:"method"`
	Template    string            `json:"template"`	
	ContentType string            `json:"content_type"`
				
	Registry  	string 		`json:"registry"`
	Repo 		string 		`json:"repo"`
	Name  		string 		`json:"name"`
	Tag 		string 		`json:"tag"`
	Ports 		[]int 		`json:"ports"`
	PortBindings 	map[string]string `json:"port_bindings"`
	Env 		[]string 	`json:"env"`
    	Links 		map[string]string `json:"links"`
	Volumes 	map[string]string `json:"volumes"`
    	PublishAllPorts bool `json:"publish_all_ports"`
}

// Auth represents a basic HTTP authentication username and password.
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
