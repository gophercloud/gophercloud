package networks

type Extension struct {
	Updated     string        `json:"updated"`
	Name        string        `json:"name"`
	Links       []interface{} `json:"links"`
	Namespace   string        `json:"namespace"`
	Alias       string        `json:"alias"`
	Description string        `json:"description"`
}
