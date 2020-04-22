package base

type ResCode int

const (
	ResCodeOk               ResCode = 1000
	ResCodeValidationError  ResCode = 2000
	ResCodeError            ResCode = 4000
	ResCodeInnerServerError ResCode = 5000
	ResCodeBizError         ResCode = 6000
)

type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
