package structs

type Request struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	PathPattern string `json:"path_pattern"`
}
