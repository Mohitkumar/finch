package httpserver

type Flow struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}
