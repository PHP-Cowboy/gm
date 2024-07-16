package errUtil

type UnmarshalErr struct {
	Msg string
}

func (e *UnmarshalErr) Error() string {
	return e.Msg
}

func NewUnmarshalErr() *UnmarshalErr {
	//包含解析失败或验证有错误
	return &UnmarshalErr{Msg: "json Verification failed"}
}
