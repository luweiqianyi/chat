package common

const (
	Success        = 0
	Failed         = 1
	ParameterError = 2
)

type ResponseHeader struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"` // omitempty: allow to omit
}

func toString(code int) string {
	switch code {
	case Success:
		return "success"
	case Failed:
		return "failed"
	case ParameterError:
		return "parameter error"
	}

	return "unknown error"
}

func NewSuccessResponse() ResponseHeader {
	return ResponseHeader{
		Result: toString(Success),
	}
}

func NewFailedResponse() ResponseHeader {
	return ResponseHeader{
		Result: toString(Failed),
	}
}

func NewParameterErrorResponse() ResponseHeader {
	return ResponseHeader{
		Result:  toString(Failed),
		Message: toString(ParameterError),
	}
}

func NewCustomizeFailedResponse(message string) ResponseHeader {
	return ResponseHeader{
		Result:  toString(Failed),
		Message: message,
	}
}
