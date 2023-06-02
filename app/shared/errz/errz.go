package errz

type ErrZ struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Option func(e *ErrZ)

func (e ErrZ) Error() string {
	return e.Msg
}

func NewErrZ(options ...Option) *ErrZ {
	errZ := &ErrZ{
		Code: Success,
		Msg:  SuccessMsg,
	}
	for _, option := range options {
		option(errZ)
	}
	return errZ
}

func WithMsg(msg string) Option {
	return func(e *ErrZ) {
		e.Msg = msg
	}
}

func WithErr(err error) Option {
	return func(e *ErrZ) {
		e.Msg = err.Error()
	}
}

func WithCode(code int64) Option {
	return func(e *ErrZ) {
		e.Code = code
	}
}

func (e ErrZ) GetCode() int64 {
	return e.Code
}
