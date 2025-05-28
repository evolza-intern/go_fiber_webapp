package models

type Service struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"` // timeout in seconds
}

type ProxyRequest struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

type ProxyResponse struct {
	StatusCode int               `json:"status_code"`
	Body       []byte            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
