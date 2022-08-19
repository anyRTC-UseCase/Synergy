package models

//用户登录参数
type ReqSignIn struct {
	//用户名
	UserName string `json:"userName"`
	//用户id(9位随机数字)
	UId string `json:"uid"`
	//用户工种
	WorkName string `json:"workName"`
	//用户类型(1:智能终端,2:专家,3:管理员)
	UserType int `json:"userType"`
}

//创建房间参数
type ReqInsertRoom struct {
	//是否加入房间(1:是,2:否)
	IsJoin int `json:"isJoin"`
}

//房间id参数
type ReqRoomId struct {
	//房间id
	RoomId string `json:"roomId"`
}

//加入房间参数
type ReqJoinRoom struct {
	//房间id
	RoomId string `json:"roomId"`
	//用户角色(1:主播,2:观众)
	UserRole int `json:"userRole"`
}

//用户离开房间参数
type ReqLeaveRoom struct {
	//房间id
	RoomId string `json:"roomId"`
	//用户id
	UId string `json:"uid"`
}

//获取分页列表参数
type ParamPageList struct {
	//页码（从1开始）
	PageNum int `json:"pageNum"`
	//单页显示条数
	PageSize int `json:"pageSize"`
}

//获取房间列表参数
type ReqGetRoomList struct {
	//页码（从1开始）
	PageNum int `json:"pageNum"`
	//单页显示条数
	PageSize int `json:"pageSize"`
	//房间状态(1:结束,2:进行中)
	RoomState int `json:"roomState"`
}

//记录用户在线心跳包信息参数
type ReqInsertUserOnlineInfo struct {
	//用户id
	UId string `json:"uid"`
	//心跳包时间戳(单位:秒)
	OptTs int `json:"optTs"`
	//用户所在房间id(用户不在房间传空)
	//RoomId string `json:"roomId"`
}

//uid
type ReqUId struct {
	//用户id
	UId string `json:"uid"`
}

//消息回调
type ReqTeamViewRtcNotify struct {
	//通知的id
	NoticeId string `json:"noticeID"`
	//开发者通知服务用于签名的secret
	Secret string `json:"secret"`
	//通知产品类型
	ProductId int `json:"productId"`
	//Rtc消息通知类型
	EventType int `json:"eventType"`
	//通知时间戳（单位：毫秒）
	NotifyMs int64 `json:"notifyMs"`
	//请求包
	Payload Payload `json:"payload"`
}

//消息回调Payload
type Payload struct {

	//频道名字
	ChannelName string `json:"channelName"`
	//频道会话Id
	ChannelSId string `json:"channelSId"`
	//应用Id
	AppId string `json:"appId"`
	//Rtc消息通知类型
	EventType int `json:"eventType"`
	//事件发生的时间戳（单位：毫秒）
	Ts int64 `json:"Ts"`
	//业务系统uId
	UId string `json:"uid"`
	//业务系统uId
	USId string `json:"uSId"`
	//平台()
	Platform int `json:"platform"`
	//用户离开的原因
	Reason int `json:"reason"`
}

//录像回调
type ReqTeamViewVodNotify struct {
	//通知的id
	NoticeId string `json:"noticeID"`
	//通知产品类型
	ProductId int `json:"productId"`
	//Rtc消息通知类型
	EventType int `json:"eventType"`
	//通知时间戳（单位：毫秒）
	NotifyMs int64 `json:"notifyMs"`
	//请求包
	Payload VodPayload `json:"payload"`
}

//录像回调Payload
type VodPayload struct {

	//频道名字
	CName string `json:"cname"`
	//业务系统uid
	UId string `json:"uid"`
	//业务系统sid
	SId string `json:"sid"`
	//消息序列号
	Sequence int `json:"sequence"`
	//事件发生的时间戳（单位：毫秒）
	SendTs int64 `json:"sendts"`
	//回调事件服务的类型(0:云端录制服务,1:录制服务,2:上传模块,4:扩展服务)
	ServiceType int `json:"serviceType"`
	//具体的消息内容(json)
	Details string `json:"details"`
}

//获取resourceId参数
type ReqGetVodResourceId struct {

	//待录制的频道名,即roomId
	Cname string `json:"cname,omitempty"`
	//字符串内容为云端录制服务在频道内使用的 UID，用于标识该录制服务，例如 "527841"。
	UID string `json:"uid,omitempty"`
	//云端录制 RESTful API 的调用时效，从成功开启云端录制并获得 sid （录制 ID）后开始计算，单位为小时。超时后，你将无法调用 query，updateLayout，和 stop 方法。resourceExpiredHour 需大于等于 1， 且小于等于 720，默认值为 72。
	ClientRequest AcquireClientRequest `json:"clientRequest,omitempty"`
}

type AcquireClientRequest struct {
	//云端录制 RESTful API 的调用时效
	ResourceExpiredHour int `json:"resourceExpiredHour"`
}

//开始录制参数
type ReqStartVodRecording struct {
	//字符串内容为云端录制服务在频道内使用的 UID，用于标识该录制服务，需要和你在 acquire 请求中输入的 UID 相同。
	UID string `json:"uid,omitempty"`
	//待录制的频道名,即roomId
	Cname         string             `json:"cname,omitempty"`
	ClientRequest StartClientRequest `json:"clientRequest,omitempty"`
}
type StartClientRequest struct {
	//用于鉴权的动态密钥。如果你的项目已启用 App 证书，则务必在该参数中传入你项目的动态秘钥。详见校验用户权限。
	Token string `json:"token"`
	//媒体流订阅的详细设置。云端录制会根据此设置订阅频道内的媒体流，并生成录制文件或截图。
	RecordingConfig RecordingConfig `json:"recordingConfig,omitempty"`
	//录制文件的详细设置。
	RecordingFileConfig RecordingFileConfig `json:"recordingFileConfig,omitempty"`
	//第三方云存储的设置。
	StorageConfig StorageConfig `json:"storageConfig,omitempty"`
}

type RecordingConfig struct {
	//最长空闲频道时间，单位为秒。如果频道内无用户的状态持续超过该时间，录制程序会自动退出。该值需大于等于 5，且小于等于 (232-1)。退出后，再次调用 start 请求，会产生新的录制文件。
	//
	//注意事项:
	//通信场景下，如果频道内有用户，但用户没有发流，不算作无用户状态。
	//直播场景下，如果频道内有观众但无主播，一旦无主播的状态超过 maxIdleTime，录制程序会自动退出。
	MaxIdleTime int `json:"maxIdleTime,omitempty"`
	//订阅的媒体流类型。
	//0：仅订阅音频
	//1：仅订阅视频
	//2：订阅音频和视频
	StreamTypes int `json:"streamTypes,omitempty"`
	//频道场景。频道场景必须与 anyRTC RTC Native/Web SDK 的设置一致，否则可能导致问题。
	//0：通信场景
	//1：直播场景
	ChannelType int `json:"channelType"`
	//视频转码的详细设置。仅适用于合流模式，单流模式下不能设置该参数。如果不设置将使用默认值。如果设置该参数，请务必填入 width、height、fps 和 bitrate 字段。请参考设置输出视频属性设置该参数。
	TranscodingConfig TranscodingConfig `json:"transcodingConfig,omitempty"`
	//预估的订阅人数峰值。在单流模式下，为必填参数。举例来说，如果 subscribeVideoUids 为 ["100","101","102"]，subscribeAudioUids 为 ["101","102","103"]，则订阅人数为 4 人。
	//
	//0：1 到 2 个 UID
	//1：3 到 7 个 UID
	//2：8 到 12 个 UID
	//3：13 到 17 个 UID
	//SubscribeUidGroup int `json:"subscribeUidGroup"` //预估的订阅人数峰值。在单流模式下，为必填参数。举例来说，如果 subscribeVideoUids 为 ["100","101","102"]，subscribeAudioUids 为 ["101","102","103"]，则订阅人数为 4 人。
	//
	//（选填）Number 类型，设置订阅的视频流类型。如果频道中有用户开启了双流模式，你可以选择订阅视频大流或者小流。
	//0：视频大流（默认），即高分辨率高码率的视频流
	//1：视频小流，即低分辨率低码率的视频流
	VideoStreamType int `json:"videoStreamType"`

	//（选填）JSONArray 类型，由 UID 组成的数组，指定订阅哪几个 UID 的音频流。数组长度不得超过 32，不推荐使用空数组。
	SubscribeAudioUids []string `json:"subscribeAudioUids"`

	//（选填）JSONArray 类型，由 UID 组成的数组，指定订阅哪几个 UID 的视频流。数组长度不得超过 32，不推荐使用空数组。
	SubscribeVideoUids []string `json:"subscribeVideoUids"`
}
type TranscodingConfig struct {
	//example: 640
	//视频的高度，单位为像素。height 不能超过 1920，且 width 和 height 的乘积不能超过 1920 * 1080，超过最大值会报错。
	Height int `json:"height,omitempty"`
	//example: 360
	//视频的宽度，单位为像素。width 不能超过 1920，且 width 和 height 的乘积不能超过 1920 * 1080，超过最大值会报错。
	Width int `json:"width,omitempty"`
	//example: 500
	//视频的码率，单位 Kbps。
	Bitrate int `json:"bitrate,omitempty"`
	//example: 15
	//视频的帧率，单位 fps。
	Fps int `json:"fps,omitempty"`
	//default: 0
	//设置视频合流布局，0、1、2 为预设的合流布局，3 为自定义合流布局。该参数设为 3 时必须设置 layoutConfig 参数。
	//
	//0: 悬浮布局。第一个加入频道的用户在屏幕上会显示为大视窗，铺满整个画布，其他用户的视频画面会显示为小视窗，从下到上水平排列，最多 4 行，每行 4 个画面，最多支持共 17 个画面。
	//1：自适应布局。根据用户的数量自动调整每个画面的大小，每个用户的画面大小一致，最多支持 17 个画面。
	//2：垂直布局。指定一个用户在屏幕左侧显示大视窗画面，其他用户的小视窗画面在右侧垂直排列，最多两列，一列 8 个画面，最多支持共 17 个画面。
	//3：自定义布局。设置 layoutConfig 参数自定义合流布局。
	MixedVideoLayout int `json:"mixedVideoLayout,omitempty"`
}
type RecordingFileConfig struct {
	//由多个字符串组成的数组，指定录制的视频文件类型。目前只支持默认值 ["hls"]，即录制生成 M3U8 和 TS 文件。
	AvFileType []string `json:"avFileType"`
}
type StorageConfig struct {
	//第三方云存储的 access key。建议提供只写权限的访问密钥（当vendor=100时，此处可为""）。
	AccessKey string `json:"accessKey,omitempty"`
	//第三方云存储指定的地区信息。 当 vendor = 0，即第三方云存储为七牛云时：
	//
	//0：Huadong
	//1：Huabei
	//2：Huanan
	//3：Beimei
	//4：Dongnanya
	//当 vendor = 1，即第三方云存储为 Amazon S3 时：
	//
	//0：US_EAST_1
	//1：US_EAST_2
	//2：US_WEST_1
	//3：US_WEST_2
	//4：EU_WEST_1
	//5：EU_WEST_2
	//6：EU_WEST_3
	//7：EU_CENTRAL_1
	//8：AP_SOUTHEAST_1
	//9：AP_SOUTHEAST_2
	//10：AP_NORTHEAST_1
	//11：AP_NORTHEAST_2
	//12：SA_EAST_1
	//13：CA_CENTRAL_1
	//14：AP_SOUTH_1
	//15：CN_NORTH_1
	//17：US_GOV_WEST_1
	//当 vendor = 2，即第三方云存储为阿里云时：
	//
	//0：CN_Hangzhou
	//1：CN_Shanghai
	//2：CN_Qingdao
	//3：CN_Beijin
	//4：CN_Zhangjiakou
	//5：CN_Huhehaote
	//6：CN_Shenzhen
	//7：CN_Hongkong
	//8：US_West_1
	//9：US_East_1
	//10：AP_Southeast_1
	//11：AP_Southeast_2
	//12：AP_Southeast_3
	//13：AP_Southeast_5
	//14：AP_Northeast_1
	//15：AP_South_1
	//16：EU_Central_1
	//17：EU_West_1
	//18：EU_East_1
	//当 vendor = 3，即第三方云存储为腾讯云时：
	//
	//0：AP_Beijing_1
	//1：AP_Beijing
	//2：AP_Shanghai
	//3：AP_Guangzhou
	//4：AP_Chengdu
	//5：AP_Chongqing
	//6：AP_Shenzhen_FSI
	//7：AP_Shanghai_FSI
	//8：AP_Beijing_FSI
	//9：AP_Hongkong
	//10：AP_Singapore
	//11：AP_Mumbai
	//12：AP_Seoul
	//13：AP_Bangkok
	//14：AP_Tokyo
	//15：NA_Siliconvalley
	//16：NA_Ashburn
	//17：NA_Toronto
	//18：EU_Frankfurt
	//19：EU_Moscow
	//当 vendor = 4，即第三方云存储为金山云时：
	//
	//0：CN_Hangzhou
	//1：CN_Shanghai
	//2：CN_Qingdao
	//3：CN_Beijing
	//4：CN_Guangzhou
	//5：CN_Hongkong
	//6：JR_Beijing
	//7：JR_Shanghai
	//8：NA_Russia_1
	//9：NA_Singapore_1
	Region int `json:"region,omitempty"`
	//第三方云存储的 bucket（当vendor=100时，bucket应为访问私有存储的接口地址）。
	Bucket string `json:"bucket,omitempty"`
	//第三方云存储的 secret key（当vendor=100时，此处可为""）。
	SecretKey string `json:"secretKey,omitempty"`
	//第三方云存储供应商。
	//
	//0：七牛云
	//1：Amazon S3
	//2：阿里云
	//3：腾讯云
	//4：金山云
	//100：私有存储
	Vendor int `json:"vendor,omitempty"`
	//由多个字符串组成的数组，指定录制文件在第三方云存储中的存储位置。举个例子，fileNamePrefix = ["directory1","directory2"]，将在录制文件名前加上前缀 "directory1/directory2/"，即 directory1/directory2/xxx.m3u8。前缀长度（包括斜杠）不得超过 128 个字符。字符串中不得出现斜杠。以下为支持的字符集范围：
	//
	//26 个小写英文字母 a-z
	//26 个大写英文字母 A-Z
	//10 个数字 0-9
	FileNamePrefix []string `json:"fileNamePrefix,omitempty"`
}

// Stop 停止云端录制
type ReqStopVodRecording struct {
	//录制的频道名。
	Cname string `json:"cname,omitempty"`
	//字符串内容为云端录制服务在频道内使用的 UID，用于标识该录制服务，需要和你在 acquire 请求中输入的 UID 相同。
	UId string `json:"uid,omitempty"`
}

type VodNotifyDetails struct {
	EventType int                 `json:"eventType"`
	MsgName   string              `json:"msgName"`
	Status    int                 `json:"status"`
	FileList  []VodNotifyFileInfo `json:"fileList"`
}

type VodNotifyFileInfo struct {
	Filename       string `json:"filename"`
	TrackType      string `json:"trackType"`
	MixedAllUser   bool   `json:"mixedAllUser"`
	UID            string `json:"uid"`
	IsPlayable     bool   `json:"isPlayable"`
	SliceStartTime int64  `json:"sliceStartTime"`
}
