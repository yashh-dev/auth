package handlers

type ResponseStatus struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
}

type Response struct {
	Content interface{}    `json:"content"`
	Status  ResponseStatus `json:"status"`
}
