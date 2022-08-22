package consts

import "time"

const (
	// IntZero 常量0
	IntZero = 0
	// IntOne 常量1
	IntOne = 1
	// IntTwo 常量2
	IntTwo = 2
	// IntThree 常量3
	IntThree = 3
	// IntFive 常量5
	IntFive = 5
	// IntTen 常量10
	IntTen = 10
	// IntThirty 常量30
	IntThirty = 30
	// IntOneZeroTwoFour 常量1024
	IntOneZeroTwoFour = 1024
)

//常量
const (
	// StrColon 冒号
	StrColon = ":"
	StrEmpty = ""
	// StrSpace 空格
	StrSpace = " "
	// StrAsterisk *
	StrAsterisk = "*"
	// StrUnderline _
	StrUnderline = "_"
	// StrAsterisk *
	StrFileSuffix = ".mp4"
)
const (
	// DftLogExpire 默认日志保存天数
	DftLogExpire = 1
	Dft1Sec      = 1 * time.Second
	Dft2Sec      = 2 * time.Second
	Dft3Sec      = 3 * time.Second
	Dft5Sec      = 5 * time.Second
	Dft10Sec     = 10 * time.Second
	Dft1Min      = time.Minute
	Dft5Min      = 5 * time.Minute
	Dft10Min     = 10 * time.Minute
	Dft15Min     = 15 * time.Minute
	Dft20Min     = 20 * time.Minute
	Dft25Min     = 25 * time.Minute
	Dft30Min     = 30 * time.Minute
	// Dft24Hour 24小时
	Dft24Hour = 24 * time.Hour
)

const (
	// RTCTokenHeader 校验header
	RTCTokenHeader = "artc-token"
	// ARAnyRtcStr anyrtc 常量
	ARAnyRtcStr = "anyrtc"
	// ARUIdStr uid 常量
	ARUIdStr = "uid"
	// ARAliYunStr 常量
	ARAliYunStr = "aliyun"
)

const (
	// ARtcChanCreate 频道创建
	ARtcChanCreate = 101
	// ARtcChanDestroy 频道销毁
	ARtcChanDestroy = 102
	// ARtcBroadcasterJoin 直播场景下主播加入频道
	ARtcBroadcasterJoin = 103
	// ARtcBroadCasterLeave 直播场景下主播离开频道
	ARtcBroadCasterLeave = 104
	// ARtcAudienceJoin 直播场景下观众进入频道
	ARtcAudienceJoin = 105
	// ARtcAudienceLeave 直播场景下观众离开频道
	ARtcAudienceLeave = 106
	// ARtcCommunicationJoin 通信场景下用户加入频道
	ARtcCommunicationJoin = 107
	// ARtcCommunicationLeave 通讯场景下用户离开频道
	ARtcCommunicationLeave = 108
	// ARtcRole2Broadcaster 直播场景下，加入频道后，用户将角色切换为主播
	ARtcRole2Broadcaster = 111
	// ARtcRole2Audience 直播场景下，加入频道后，用户将角色切换为观众
	ARtcRole2Audience = 112
)

//默认的录制参数
const (
	// DftResExpire resourceId的默认缓存时间
	DftResExpire = 72
	// MaxResExpire resourceId的最大时间不能超过720小时
	MaxResExpire = 720
	// DftVodModMix 混流模式
	DftVodModMix           = "mix"
	DftVodMaxIdleTime      = 60
	DftVodChannelType      = 1
	DftVodAudioProfile     = 0
	DftVodSnapInterval     = 10
	DftVodVideoStreamType  = 0
	DftVodStreamType       = 2
	DftVodTransWidth       = 360
	DftVodTransHeight      = 640
	DftVodTransFps         = 15
	DftVodTransBitrate     = 500
	DftVodTransVideoLayout = 1
	DftSubscribeUidGroup   = 0
	DftVodBgColor          = "#000000"
	DftVodContentType      = "application/json;charset=utf-8"
)

//录制媒体类型
const (
	VodStreamHls = "hls"
	VodStreamMp4 = "mp4"
)

//云端录制回调事件
const (
	// 录制文件已上传至指定的第三方云存储
	ARtcVodUploaded = 31
	// 录制服务已启动
	ARtcVodRecorderStarted = 40
	// 录制组件已退出
	ARtcVodRecorderLeave = 41
)

// 常量
const (
	//用户离开时间
	UserLeaveTime = 0
	//房间名称的后缀
	RoomNameSuffix = "的频道"
)

// 房间状态(1:结束,2:进行中,3:转码中)
const (
	//1:录像结束
	RoomStateClosed = 1
	//2:进行中
	RoomStateOpen = 2
	//3:转码中
	RoomStateMixed = 3
	//4:通话结束
	RoomStateCallFinished = 4
)

// 是否加入房间(1:是,2:否)
const (
	//1:是
	IsJoinRoom = 1
	//2:否
	IsNotJoinRoom = 2
)

//用户类型(1:智能终端,2:专家,3:管理员)
const (
	//1:智能终端
	UserTypeTerminal = 1
	//2:专家
	UserTypeSpecialist = 2
	//3:管理员
	UserTypeAdmin = 3
)

//用户角色(1:主播,2:观众)
const (
	// 1:主播
	UserRoleHost = 1
	// 2:观众
	UserRoleAudience = 2
)

//专家状态(1:通话中,2:空闲,3:离线)
const (
	// 1:通话中
	UserStateBusy = 1
	// 2:空闲
	UserStateLeisure = 2
	// 3:离线
	UserStateOffline = 3
)

//0：七牛云
//1：Amazon S3
//2：阿里云
//3：腾讯云
//4：金山云
//100：私有存储
const (
	// VendorAliYun 2：阿里云
	VendorAliYun = 2
)
