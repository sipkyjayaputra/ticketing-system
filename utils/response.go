package utils

type ResponseContainer struct {
	Response Response `json:"response"`
}

type Response struct {
	StatusCode      int           `json:"status_code"`
	Success         bool          `json:"success"`
	ResponseMessage *string       `json:"response_message"`
	Errors          []string      `json:"errors"`
	Data            interface{}   `json:"data"`
	Info            *ResponseInfo `json:"info,omitempty"`
}

type ErrorContainer struct {
	Response Response `json:"response"`
}

type ResponseInfo struct {
	Limit    int `json:"limit"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}
