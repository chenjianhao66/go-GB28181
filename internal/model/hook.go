package model

const (
	RespondSuccess     = 0    //执行成功
	RespondException   = -400 //代码抛异常
	RespondInvalidArgs = -300 //参数不合法
	RespondSqlFailed   = -200 //sql执行失败
	RespondAuthFailed  = -100 //鉴权失败
	RespondOtherFailed = -1   //业务代码执行失败，
)

const (
	ParseParamFail = iota + 1
)

const (
	SuccessMsg        = "success"
	ParseParamFailMsg = "parser on_play_hook interface param fail, auth fail and not allow play"
)

// HookReply hook事件默认回复
type HookReply struct {
	// 0代表允许，其他均为不允许
	Code int `json:"code,omitempty"`

	// 当code不为0时，msg字段应给出相应提示
	Msg string `json:"msg,omitempty"`
}

// 流媒体服务心跳
type (
	ServerKeepalive struct {
		MediaServerId string              `json:"mediaServerId,omitempty"`
		Data          ServerKeepaliveData `json:"data"`
	}

	ServerKeepaliveData struct {
		Buffer                int `json:"Buffer,omitempty"`
		BufferLikeString      int `json:"BufferLikeString,omitempty"`
		BufferList            int `json:"BufferList,omitempty"`
		BufferRaw             int `json:"BufferRaw,omitempty"`
		Frame                 int `json:"Frame,omitempty"`
		FrameImp              int `json:"FrameImp,omitempty"`
		MediaSource           int `json:"MediaSource,omitempty"`
		MultiMediaSourceMuxer int `json:"MultiMediaSourceMuxer,omitempty"`
		RtmpPacket            int `json:"RtmpPacket,omitempty"`
		RtpPacket             int `json:"RtpPacket,omitempty"`
		Socket                int `json:"Socket,omitempty"`
		TcpClient             int `json:"TcpClient,omitempty"`
		TcpServer             int `json:"TcpServer,omitempty"`
		TcpSession            int `json:"TcpSession,omitempty"`
		UdpServer             int `json:"UdpServer,omitempty"`
		UdpSession            int `json:"UdpSession,omitempty"`
	}
)

// OnPlayHookParam 播放器鉴权hook事件
type OnPlayHookParam struct {
	// 流应用名
	App string `json:"app,omitempty" `

	// TCP链接唯一ID
	Id string `json:"id,omitempty"`

	// 播放器ip
	Ip string `json:"ip,omitempty"`

	// 播放url参数
	Params string `json:"params,omitempty"`

	// 播放器端口号
	Port int `json:"port,omitempty"`

	// 播放的协议，可能是rtsp、rtmp、http
	Schema string `json:"schema,omitempty"`

	// 流ID
	Stream string `json:"stream,omitempty"`

	// 流虚拟主机
	Vhost string `json:"vhost,omitempty"`

	// 服务器id,通过配置文件设置
	MediaServerId string `json:"mediaServerId,omitempty"`
}

type (
	// OnPublishHookParam rtsp、rtmp、rtp推流鉴权事件
	OnPublishHookParam struct {
		App           string `json:"app,omitempty"`
		Id            string `json:"id,omitempty"`
		Ip            string `json:"ip,omitempty"`
		Params        string `json:"params,omitempty"`
		Port          int    `json:"port,omitempty"`
		Schema        string `json:"schema,omitempty"`
		Stream        string `json:"stream,omitempty"`
		Vhost         string `json:"vhost,omitempty"`
		MediaServerId string `json:"mediaServerId,omitempty"`
	}

	OnPublishHookReply struct {
		HookReply

		// 是否转换成hls协议
		EnableHls bool `json:"enable_hls,omitempty"`

		// 是否允许mp4录制
		EnableMp4 bool `json:"enable_mp4,omitempty"`

		// 是否转rtsp协议
		EnableRtsp bool `json:"enable_rtsp,omitempty"`

		// 是否转rtmp/flv协议
		EnableRtmp bool `json:"enable_rtmp,omitempty"`

		// 是否转http-ts/ws-ts协议
		EnableTs bool `json:"enable_ts,omitempty"`

		// 是否转http-fmp4/ws-fmp4协议
		EnableFmp4 bool `json:"enable_fmp4,omitempty"`

		// 转协议时是否开启音频
		EnableAudio bool `json:"enable_audio,omitempty"`

		// 转协议时，无音频是否添加静音aac音频
		AddMuteAudio bool `json:"add_mute_audio,omitempty"`

		// mp4录制文件保存根目录，置空使用默认
		Mp4SavePath string `json:"mp4_save_path,omitempty"`

		// mp4录制切片大小，单位秒
		Mp4MaxSecond int `json:"mp4_max_second,omitempty"`

		// hls文件保存保存根目录，置空使用默认
		HlsSavePath string `json:"hls_save_path,omitempty"`

		// 	断连续推延时，单位毫秒，置空使用配置文件默认值
		ContinuePushMs uint32 `json:"continue_push_ms,omitempty"`

		// MP4录制是否当作观看者参与播放人数计数
		Mp4AsPlayer bool `json:"mp4_as_player,omitempty"`

		// 该流是否开启时间戳覆盖
		ModifyStamp bool `json:"modify_stamp,omitempty"`
	}
)

func NewOnPublishDefaultReply() OnPublishHookReply {
	return OnPublishHookReply{
		HookReply: HookReply{
			Code: RespondSuccess,
			Msg:  SuccessMsg,
		},
		AddMuteAudio:   true,
		ContinuePushMs: 10000,
		EnableAudio:    true,
		EnableFmp4:     true,
		EnableHls:      true,
		EnableMp4:      false,
		EnableRtmp:     true,
		EnableRtsp:     true,
		EnableTs:       true,
		HlsSavePath:    "/hls_save_path/",
		ModifyStamp:    false,
		Mp4AsPlayer:    false,
		Mp4MaxSecond:   3600,
		Mp4SavePath:    "/mp4_save_path/",
	}
}

// 流改变事件hook参数
type (
	OnStreamChangedParam struct {
		Register      bool       `json:"regist,omitempty"`
		AliveSecond   uint       `json:"aliveSecond,omitempty"`
		App           string     `json:"app,omitempty"`
		BytesSpeed    uint       `json:"bytesSpeed,omitempty"`
		CreateStamp   uint       `json:"createStamp,omitempty"`
		MediaServerId string     `json:"mediaServerId,omitempty"`
		OriginSock    OriginData `json:"originSock"`

		// 产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7,rtc_push=8
		OriginType int `json:"originType,omitempty"`

		// 源类型名称
		OriginTypeStr string `json:"originTypeStr,omitempty"`
		OriginUrl     string `json:"originUrl,omitempty"`

		// 协议观看人数
		ReaderCount int `json:"readerCount,omitempty"`

		Schema           string  `json:"schema,omitempty"`
		Stream           string  `json:"stream,omitempty"`
		TotalReaderCount int     `json:"totalReaderCount,omitempty"`
		Tracks           []Track `json:"tracks,omitempty"`
		Vhost            string  `json:"vhost,omitempty"`
	}

	// OriginData 对端数据
	OriginData struct {
		Identifier string `json:"identifier,omitempty"`
		LocalIp    string `json:"local_ip,omitempty"`
		LocalPort  string `json:"local_port,omitempty"`
		PeerIp     string `json:"peer_ip,omitempty"`
		PeerPort   string `json:"peer_port,omitempty"`
	}

	// Track 音视频轨道
	Track struct {
		// 音频通道数
		Channels int `json:"channels,omitempty"`

		// 编解码id，H264 = 0, H265 = 1, AAC = 2, G711A = 3, G711U = 4
		CodecId int `json:"codec_id,omitempty"`

		// 编码类型名称
		CodecIdName string `json:"codec_id_name,omitempty"`

		// 编码类型，video=0，audio=1
		CodecType int `json:"codec_type,omitempty"`

		// 轨道是否准备就绪
		Ready bool `json:"ready,omitempty"`

		// 音频采样位数
		SampleBit int `json:"sample_bit,omitempty"`

		// 音频采样率
		SampleRate int `json:"sample_rate,omitempty"`
	}
)

type (
	// OnStreamNoneReader 流无人观看时事件参数
	OnStreamNoneReader struct {
		MediaServerId string `json:"mediaServerId,omitempty"`
		App           string `json:"app,omitempty"`
		Schema        string `json:"schema,omitempty"`
		Stream        string `json:"stream,omitempty"`
		Vhost         string `json:"vhost,omitempty"`
	}

	// OnStreamNoneReaderReply 流无人观看事件回复
	OnStreamNoneReaderReply struct {
		// 固定返回0
		Code int `json:"code,omitempty"`
		// 是否关闭该流，包括推流和拉流
		Close bool `json:"close,omitempty"`
	}
)

type (
	MediaConfig struct {
		RemoteIp string
		// 是否调试http api,启用调试后，会打印每次http请求的内容和回复
		ApiDebug string `json:"api.apiDebug,omitempty" mapstructure:"api.apiDebug"`
		// 默认截图图片
		ApiDefaultSnap string `json:"api.defaultSnap,omitempty" mapstructure:"api.defaultSnap"`
		// 截图保存路径根目录
		ApiSnapRoot string `json:"api.snapRoot,omitempty" mapstructure:"api.snapRoot"`
		// 密钥
		ApiSecret string `json:"api.secret,omitempty" mapstructure:"api.secret"`

		// ffmpeg相关
		FfmpegBin        string `json:"ffmpeg.bin,omitempty"`
		FfmpegCmd        string `json:"ffmpeg.cmd,omitempty"`
		FfmpegSnap       string `json:"ffmpeg.snap,omitempty"`
		FfmpegLog        string `json:"ffmpeg.log,omitempty"`
		FfmpegRestartSec string `json:"ffmpeg.restart_sec,omitempty"`

		// protocol协议相关
		ProtocolModifyStamp    string `json:"protocol.modify_stamp,omitempty"`
		ProtocolEnableAudio    string `json:"protocol.enable_audio,omitempty"`
		ProtocolAddMuteAudio   string `json:"protocol.add_mute_audio,omitempty"`
		ProtocolContinuePushMs string `json:"protocol.continue_push_ms,omitempty"`
		ProtocolEnableHls      string `json:"protocol.enable_hls,omitempty"`
		ProtocolEnableMp4      string `json:"protocol.enable_mp4,omitempty"`
		ProtocolEnableRtsp     string `json:"protocol.enable_rtsp,omitempty"`
		ProtocolEnableRtmp     string `json:"protocol.enable_rtmp,omitempty"`
		ProtocolEnableTs       string `json:"protocol.enable_ts,omitempty"`
		ProtocolEnableFmp4     string `json:"protocol.enable_fmp4,omitempty"`
		ProtocolMp4AsPlayer    string `json:"protocol.mp4_as_player,omitempty"`
		ProtocolMp4MaxSecond   string `json:"protocol.mp4_max_second,omitempty"`
		ProtocolMp4SavePath    string `json:"protocol.mp4_save_path,omitempty"`
		ProtocolHlsSavePath    string `json:"protocol.hls_save_path,omitempty"`
		ProtocolHlsDemand      string `json:"protocol.hls_demand,omitempty"`
		ProtocolRtspDemand     string `json:"protocol.rtsp_demand,omitempty"`
		ProtocolRtmpDemand     string `json:"protocol.rtmp_demand,omitempty"`
		ProtocolTsDemand       string `json:"protocol.ts_demand,omitempty"`
		ProtocolFmp4Demand     string `json:"protocol.fmp4_demand,omitempty"`

		// 通用配置
		GeneralEnableVhost             string `json:"general.enableVhost,omitempty"`
		GeneralFlowThresHold           string `json:"general.flowThreshold,omitempty"`
		GeneralMaxWaitMs               string `json:"general.maxWaitMS,omitempty"`
		GeneralStreamNoneReaderDelayMs string `json:"general.streamNoneReaderDelayMs,omitempty"`
		GeneralResetWhenReplay         string `json:"general.resetWhenReplay,omitempty"`
		GeneralMergeWriteMs            string `json:"general.mergeWriteMs,omitempty"`
		GeneralMediaServerId           string `json:"general.mediaServerId,omitempty"`

		// hls配置
		HlsFileBufSize       string `json:"hls.fileBufSize,omitempty"`
		HlsSegDur            string `json:"hls.segDur,omitempty"`
		HlsSegRetain         string `json:"hls.segRetain,omitempty"`
		HlsSegKeep           string `json:"hls.segKeep,omitempty"`
		HlsBroadcastRecordTs string `json:"hls.broadcastRecordTs,omitempty"`
		HlsDeleteDelaySec    string `json:"hls.deleteDelaySec,omitempty"`

		// hook配置
		HookAdminParams        string `json:"hook.admin_params,omitempty"`
		HookEnable             string `json:"hook.enable,omitempty"`
		HookOnFlowReport       string `json:"hook.on_flow_report,omitempty"`
		HookOnHttpAccess       string `json:"hook.on_http_access,omitempty"`
		HookOnPlay             string `json:"hook.on_play,omitempty"`
		HookOnPublish          string `json:"hook.on_publish,omitempty"`
		HookOnRecordMp4        string `json:"hook.on_record_mp4,omitempty"`
		HookOnRecordTs         string `json:"hook.on_record_ts,omitempty"`
		HookOnRtspAuth         string `json:"hook.on_rtsp_auth,omitempty"`
		HookOnRtspRealm        string `json:"hook.on_rtsp_realm,omitempty"`
		HookOnShellLogin       string `json:"hook.on_shell_login,omitempty"`
		HookOnStreamChanged    string `json:"hook.on_stream_changed,omitempty"`
		HookOnStreamNoneReader string `json:"hook.on_stream_none_reader,omitempty"`
		HookOnStreamNotFound   string `json:"hook.on_stream_not_found,omitempty"`
		HookOnServerStared     string `json:"hook.on_server_stared,omitempty"`
		HookOnServerKeepalive  string `json:"hook.on_server_keepalive,omitempty"`
		HookOnSendRtpStopped   string `json:"hook.on_send_rtp_stopped,omitempty"`
		HookOnRtpServerTimeout string `json:"hook.on_rtp_server_timeout,omitempty"`
		HookTimeoutSec         string `json:"hook.timeoutSec,omitempty"`
		HookAliveInterval      string `json:"hook.alive_interval,omitempty"`
		HookRetry              string `json:"hook.retry,omitempty"`
		HookRetryDelay         string `json:"hook.retry_delay,omitempty"`

		// http配置
		HttpCharSet           string `json:"http.charSet,omitempty"`
		HttpKeepAliveSecond   string `json:"http.keepAliveSecond,omitempty"`
		HttpMaxReqSize        string `json:"http.maxReqSize,omitempty"`
		HttpPort              string `json:"http.port,omitempty"`
		HttpRootPath          string `json:"http.rootPath,omitempty"`
		HttpSendBufSize       string `json:"http.sendBufSize,omitempty"`
		HttpSSLPort           string `json:"http.sslport,omitempty"`
		HttpDirMenu           string `json:"http.dirMenu,omitempty"`
		HttpVirtualPath       string `json:"http.virtualPath,omitempty"`
		HttpForbidCacheSuffix string `json:"http.forbidCacheSuffix,omitempty"`

		//record 录像相关配置
		RecordAppName     string `json:"record.appName,omitempty"`
		RecordFileBufSize string `json:"record.fileBufSize,omitempty"`
		RecordSampleMs    string `json:"record.sampleMS,omitempty"`
		RecordFastStart   string `json:"record.fastStart,omitempty"`
		RecordFileRepeat  string `json:"record.fileRepeat,omitempty"`

		// rtmp
		RtmpHandshakeSecond string `json:"rtmp.handshakeSecond,omitempty"`
		RtmpKeepAliveSecond string `json:"rtmp.keepAliveSecond,omitempty"`
		RtmpModifyStamp     string `json:"rtmp.modifyStamp,omitempty"`
		RtmpPort            string `json:"rtmp.port,omitempty"`
		RtmpSSLPort         string `json:"rtmp.sslport,omitempty"`

		// rtp
		RtpAudioMtuSize string `json:"rtp.audioMtuSize,omitempty"`
		RtpVideoMtuSize string `json:"rtp.videoMtuSize,omitempty"`
		RtpMaxSize      string `json:"rtp.maxSize,omitempty"`
		RtpLowLatency   string `json:"rtp.lowLatency,omitempty"`

		// rtp proxy
		RtpProxyDumpDir    string `json:"rtp_proxy.dumpDir,omitempty"`
		RtpProxyPort       string `json:"rtp_proxy.port,omitempty"`
		RtpProxyTimeoutSec string `json:"rtp_proxy.timeoutSec,omitempty"`
		RtpProxyPortRange  string `json:"rtp_proxy.port_range,omitempty"`
		RtpProxyH264Pt     string `json:"rtp_proxy.h264_pt,omitempty"`
		RtpProxyH265Pt     string `json:"rtp_proxy.h265_pt,omitempty"`
		RtpProxyPsPt       string `json:"rtp_proxy.ps_pt,omitempty"`
		RtpProxyOpusPt     string `json:"rtp_proxy.opus_pt,omitempty"`

		// rtsp
		RtspAuthBasic       string `json:"rtsp.authBasic,omitempty"`
		RtspDirectProxy     string `json:"rtsp.directProxy,omitempty"`
		RtspHandshakeSecond string `json:"rtsp.handshakeSecond,omitempty"`
		RtspKeepAliveSecond string `json:"rtsp.keepAliveSecond,omitempty"`
		RtspPort            string `json:"rtsp.port,omitempty"`
		RtspSSLPort         string `json:"rtsp.sslport,omitempty"`
		RtspLowLatency      string `json:"rtsp.lowLatency,omitempty"`
	}
)
