package httpserver

type Flow struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}
