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
	codeGetUserId
)

var (
	ErrInvalidParam = NewErrZ(WithCode(CodeInvalidParam), WithMsg("invalid param"))

	ErrGetUserId    = NewErrZ(WithCode(codeGetUserId), WithMsg("get user id failed"))
	ErrNoPermission = NewErrZ(WithCode(CodeNoPermission), WithMsg("no permission"))
	ErrRpcCall      = NewErrZ(WithCode(CodeRpcCall), WithMsg("call rpc service failed"))
	ErrServiceBusy  = NewErrZ(WithCode(CodeServiceBusy), WithMsg("service busy"))
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

const (
	codeFileService int64 = 30000 + iota

	codeGenerateUrl
	codeGetFileList
)

var (
	ErrFileService = NewErrZ(WithCode(codeFileService), WithMsg("file service busy"))

	ErrGenerateUrl = NewErrZ(WithCode(codeGenerateUrl), WithMsg("generate url failed"))
	ErrGetFileList = NewErrZ(WithCode(codeGetFileList), WithMsg("get file list failed"))
)
