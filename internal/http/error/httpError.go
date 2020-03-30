package httperror

type HTTPError struct {
	Cause error        `json:"-"`
	Info  ErrorMessage `json:"message"`
	Code  int          `json:"code"`
}

type ErrorMessage struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
