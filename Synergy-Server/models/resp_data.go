package models

//---------------------file------------------------
type FileItem struct {
	//文件url
	Url string `json:"url"`
	//文件名字
	FileName string `json:"name"`
}

//用户信息
type UserInfo struct {
	//用户名
	UserName string `json:"userName"`
	//用户id
	UId string `json:"uid"`
	//用户工种
	WorkName string `json:"workName"`
	//用户类型(1:智能终端,2:专家,3:管理员)
	UserType int `json:"userType"`
	//用户登录时间戳(单位:秒)
	UserTs int `json:"userTs"`
}

//用户登录响应信息
type RespSignIn struct {
	//用户名
	UserName string `json:"userName"`
	//用户id
	UId string `json:"uid"`
	//用户工种
	WorkName string `json:"workName"`
	//用户类型(1:智能终端,2:专家,3:管理员)
	UserType int `json:"userType"`
	//用户登录时间戳(单位:秒)
	UserTs int `json:"userTs"`

	//appId
	AppId string `json:"appId"`
	//rtmToken
	RtmToken string `json:"rtmToken"`
}

//创建房间响应信息
type RespInsertRoom struct {
	//房间id
	RoomId string `json:"roomId"`
	//房主uid
	RoomHostId string `json:"roomHostId"`
	//房间名称
	RoomName string `json:"roomName"`
	//房间状态(1:结束,2:进行中,3:转码中)
	RoomState int `json:"roomState"`
	//房间创建时间戳(单位:秒)
	RoomTs int `json:"roomTs"`
	//rtcToken(不加入房间时为空)
	RtcToken string `json:"rtcToken"`
}

//房间信息
type RoomInfo struct {
	//房间id
	RoomId string `json:"roomId"`
	//房间名称
	RoomName string `json:"roomName"`
	//房主uid
	RoomHostId string `json:"roomHostId"`
	//房间状态(1:结束,2:进行中,3:转码中)
	RoomState int `json:"roomState"`
	//房间创建时间戳(单位:秒)
	RoomTs int `json:"roomTs"`
	//房间录制在频道内使用的 UID
	RoomVodUId string `json:"roomVodUId"`
	//房间录像的resource ID
	RoomVodResourceId string `json:"roomVodResourceId"`
	//房间录像的录制 ID
	RoomVodSId string `json:"roomVodSId"`
}

//进行中房间信息
type OnGoingRoomInfo struct {
	//房间id
	RoomId string `json:"roomId"`
	//房间名称
	RoomName string `json:"roomName"`
	//房主uid
	RoomHostId string `json:"roomHostId"`
	//房间状态(1:结束,2:进行中,3:转码中)
	RoomState int `json:"roomState"`
	//房间创建时间戳(单位:秒)
	RoomTs int `json:"roomTs"`
	//房主名
	UserName string `json:"userName"`
}

//获取进行中房间列表响应信息
type RespOngoingRoomList struct {
	//房间列表
	List []OnGoingRoomInfo `json:"list"`
	//总数量
	TotalNum int `json:"totalNum"`
}

//结束房间信息
type FinishedRoomInfo struct {
	//房间id
	RoomId string `json:"roomId"`
	//房间名称
	RoomName string `json:"roomName"`
	//房主uid
	RoomHostId string `json:"roomHostId"`
	//房间状态(1:结束,2:进行中,3:转码中)
	RoomState int `json:"roomState"`
	//房间创建时间戳(单位:秒)
	RoomTs int `json:"roomTs"`
	//房间录制开始时间(单位:毫秒)
	RoomStartTs int `json:"roomStartTs"`
	//房间录制结束时间(单位:毫秒)
	RoomStopTs int `json:"roomStopTs"`
	//房间录像文件url
	RoomFileUrl string `json:"roomFileUrl"`
	//房主名
	UserName string `json:"userName"`
}

//结束房间列表响应信息
type RespFinishedRoomList struct {
	//房间列表
	List []FinishedRoomInfo `json:"list"`
	//总数量
	TotalNum int `json:"totalNum"`
}

//加入房间响应信息
type RespJoinRoom struct {
	//房间信息
	RoomInfo RoomInfo `json:"roomInfo"`
	//rtcToken
	RtcToken string `json:"rtcToken"`
}

//获取专家列表响应信息
type RespGetSpecialist struct {
	//专家信息列表
	List []SpecialistUserInfo `json:"list"`
	//总数量
	TotalNum int `json:"totalNum"`
}

//专家信息
type SpecialistUserInfo struct {
	//用户名
	UserName string `json:"userName"`
	//用户id
	UId string `json:"uid"`
	//用户工种
	WorkName string `json:"workName"`
	//用户类型(1:智能终端,2:专家,3:管理员)
	UserType int `json:"userType"`
	//用户登录时间戳(单位:秒)
	UserTs int `json:"userTs"`
	//用户状态(1:通话中,2:空闲,3:离线)
	UserState int `json:"userState"`
	//专家所在房间id
	RoomId string `json:"roomId"`
}

//用户加入房间信息
type UserJoinInfo struct {
	//房间id
	RoomId string `json:"roomId"`
	//用户id
	UId string `json:"uid"`
	//用户角色(1:主播,2:观众)
	UserRole int `json:"userRole"`
	//用户进入房间的时间戳（单位:秒）
	JoinTime int `json:"joinTime"`
	//用户离开房间的时间戳（单位:秒）
	LeaveTime int `json:"leaveTime"`
}

//云端录制资源响应信息
type AcquireResponse struct {
	//响应码
	Code int `json:"Code"`
	//云端录制资源 resource ID，使用这个 resource ID 可以开始一段云端录制。这个 resource ID 的有效期为 5 分钟，超时需要重新请求。
	ResourceId string `json:"resourceId,omitempty"`
}

//云端录制资源Body响应信息
type AcquireBody struct {
	//云端录制资源 resource ID，使用这个 resource ID 可以开始一段云端录制。这个 resource ID 的有效期为 5 分钟，超时需要重新请求。
	ResourceId string `json:"resourceId,omitempty"`
}

//获取resourceId响应信息
type RespGetVodResourceId struct {
	//响应码
	Code int `json:"Code"`
	//云端录制资源 resource ID，使用这个 resource ID 可以开始一段云端录制。这个 resource ID 的有效期为 5 分钟，超时需要重新请求。
	Body *AcquireBody `json:"Body,omitempty"`
}

//开始录制响应body
type StartBody struct {
	//云端录制使用的 resource ID。
	ResourceId string `json:"resourceId,omitempty"`
	//录制 ID。成功开始云端录制后，你会得到一个 sid （录制 ID)。该 ID 是一次录制周期的唯一标识。
	Sid string `json:"sid,omitempty"`
}

//开始录制响应信息
type RespStartVodRecording struct {
	//响应码
	Code int `json:"Code"`
	//响应体
	Body *StartBody `json:"Body,omitempty"`
}

//
type StartResp struct {
	//云端录制使用的 resource ID。
	ResourceId string `json:"resourceId,omitempty"`
	//会话ID
	SessionId string `json:"sid,omitempty"`
	//响应码
	Code int `json:"Code"`
}

//停止录像响应返回
type RespStopResponse struct {
	//响应码
	Code int `json:"Code"`
	//响应体
	Body StopResponse `json:"Body"`
}

//停止云端录制响应信息
type StopResponse struct {
	//云端录制使用的 resource ID。
	ResourceId string `json:"resourceId,omitempty"`
	//录制 ID。成功开始云端录制后，你会得到一个 sid （录制 ID)。该 ID 是一次录制周期的唯一标识。
	Sid string `json:"sid,omitempty"`
	//服务器返回的具体信息。
	ServerResponse StopServerResponse `json:"serverResponse,omitempty"`
}

//停止云端录制服务器返回的具体信息
type StopServerResponse struct {
	//表示 fileList 字段的数据格式。如果你设置了 snapshotConfig，则不会返回该字段。
	//
	//"string"：fileList 为 String 类型。合流模式下，fileListMode 为 "string"。
	//"json"：fileList 为 JSONArray 类型。单流模式下，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode,omitempty"`
	//当 fileListMode 为 "string" 时，fileList 为 String 类型，录制产生的 M3U8 文件的文件名。当 fileListMode 为 "json" 时, fileList 为 JSONArray 类型，由每个录制文件的具体信息组成的数组。如果你设置了 snapshotConfig，则不会返回该字段。
	FileList []FileInfo `json:"fileList,omitempty"`
	//当前录制上传的状态。
	//
	//"uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
	//"backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了 anyRTC 云备份。anyRTC 服务器会自动将这部分文件继续上传至指定的第三方云存储。
	//"unknown"：未知状态。
	UploadingStatus string `json:"uploadingStatus,omitempty"`
}

//录制文件的具体信息
type FileInfo struct {
	//String 类型，录制产生的 M3U8 文件和 MP4 文件的文件名
	FileName string `json:"filename"`
	//Boolean 类型，是否可以在线播放。
	//true：可以在线播放。
	//false：无法在线播放。
	IsPlayable bool `json:"isPlayable"`
	//Boolean 类型，用户是否是分开录制的。
	//true：所有用户合并在一个录制文件中。
	//false：每个用户分开录制。
	MixedAllUser bool `json:"mixedAllUser"`
	//Number 类型，该文件的录制开始时间，Unix 时间戳，单位为毫秒。
	SliceStartTime int64 `json:"sliceStartTime"`
	//String 类型，录制文件的类型。
	//"audio"：纯音频文件。
	//"video"：纯视频文件。
	//"audio_and_video"：音视频文件。
	TrackType string `json:"trackType"`
	//String 类型，用户 UID，表示录制的是哪个用户的音频流或视频流。合流录制模式下，uid 为 "0"
	UID string `json:"uid"`
}
