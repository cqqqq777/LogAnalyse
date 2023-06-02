package errz

const (
	Success    int64 = 0
	SuccessMsg       = "success"

	FailMsg = "fail"
)

const (
	CodeInvalidParam int64 = 10000 + iota
	CodeTokenInvalid
	CodeNoPermission
	CodeRpcCall
	CodeServiceBusy
)

var (
	ErrInvalidParam = NewErrZ(WithCode(CodeInvalidParam), WithMsg("invalid param"))
)

const (
	codeUserService int64 = 20000 + iota

	codeGenerateID
	codeGenerateToken
	codeUserExist
	codeUserNotFound
	codeWrongPassword
)

var (
	ErrUserService = NewErrZ(WithCode(codeUserService), WithMsg("user service busy"))

	ErrGenerateID    = NewErrZ(WithCode(codeGenerateID), WithMsg("generate user id failed"))
	ErrGenerateToken = NewErrZ(WithCode(codeGenerateToken), WithMsg("generate token failed"))
	ErrUserExist     = NewErrZ(WithCode(codeUserExist), WithMsg("user exist"))
	ErrUserNotFound  = NewErrZ(WithCode(codeUserNotFound), WithMsg("user not found"))
	ErrWrongPassword = NewErrZ(WithCode(codeWrongPassword), WithMsg("wrong password"))
)
