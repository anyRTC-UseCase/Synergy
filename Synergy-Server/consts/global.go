package consts

//response code.
const (
	ErrCodeOk             = 0
	ErrCodeRedis          = -1
	ErrCodeMysql          = -2
	ErrCodeApiMethod      = -3
	ErrCodeParamMissing   = -4
	ErrCodeHttpErr        = -5
	ErrCodeDbError        = -6
	ErrCodeParamInvalid   = -7
	ErrCodeInvalidRequest = 400
	ErrCodeUnauthorized   = 401
	ErrCodeInternal       = 500
	ErrCodeFailed         = 800
	//用户不存在
	ErrCodeUserNotExists = 1000
	//房间不存在
	ErrCodeRoomNotExists = 1001
	//生成token失败
	ErrCodeBuildTokenErr = 1002
	//登录的用户类型和uid对应的数据库用户类型不相同
	ErrCodeUserTypeDifferent = 1003
	//用户已在其他地方登陆
	ErrCodeUserAlreadyLogin = 1100
)

//response message.
const (
	ErrOk                  = "success."
	ErrFailed              = "failed."
	ErrJsonMarshal         = "internal, json.Marshal"
	ErrJsonUnmarshal       = "internal, json.Unmarshal"
	ErrDbError             = "internal, database error"
	ErrDbQueryError        = "sql query error"
	ErrAdminNotExists      = "admin not exist"
	ErrParamInvalid        = "invalid parameters"
	ErrNoPermissionDisplay = "you have no permission to display"
	ErrParamMissing        = "params missing"
	ErrInvalidRequest      = "request is invalid"

	ErrFileNotExists  = "file not exists."
	ErrFileReadFailed = "file read failed."
	//用户不存在
	ErrUserNotExists = "user not exists."
	//房间不存在
	ErrRoomNotExists = "room not exists."
	//生成token失败
	ErrBuildTokenErr = "build token error."
	//登录的用户类型和uid对应的用户类型不相同
	ErrUserTypeDifferent = "login database uid-usertype different"
)
