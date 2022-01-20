package controllers

// HTTPResponse 공통으로 사용할 응답 메시지
type HTTPResponse struct {
	Result  string      `yaml:"result" json:"result"`
	Code    int         `yaml:"code" json:"code"`
	Message string      `yaml:"message" json:"message"`
	Data    interface{} `yaml:"data" json:"data"`
}

// For Result
const (
	STRSuccess string = "success"
	STRError   string = "error"
)

//const (
//	HTTPErrCode1000 int = 1000 + iota
//	HTTPErrCode1001
//)
//
//// For Message
//var (
//	HTTPErrMsg map[int]string
//)

// ReturnError  ...
func (HTTPResponse) ReturnError(code int, msg string) *HTTPResponse {
	return &HTTPResponse{
		Result:  STRError,
		Code:    code,
		Message: msg,
	}
}

// ReturnSuccess ...
func (HTTPResponse) ReturnSuccess(result interface{}) *HTTPResponse {
	return &HTTPResponse{
		Result: STRSuccess,
		Data:   result,
	}
}
