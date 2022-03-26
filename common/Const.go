package common

/*
 * @Description: Do not edit
 * @Author: Jianxuesong
 * @Date: 2021-06-11 17:16:09
 * @LastEditors: Jianxuesong
 * @LastEditTime: 2021-06-11 17:16:41
 * @FilePath: /Melon/app/common/Const.go
 */

const (
	// STATUS_OK 成功
	STATUS_OK      = 0    // 成功
	STATUS_NO_DATA = 2    // 无数据
	ERROR_PARAM    = 1000 // 参数错误
	ERROR_INTER    = 1001 // 服务异常|内部错误
	ERROR_CAPTCHA  = 1002 // 验证码错误
	ERROR_AUTH     = 1003 // 验证码错误

	ERROR_LOGIN = 2001

	UPLOAD_SUCCESS = 3001
	UPLOAD_FAILED  = 3002

	ERROR_REQUEST_EXPIRED      = 1005      //请求已过期
	ERROR_IP_BLOCKED           = 1006      //IP被屏蔽
	ERROR_TK                   = 1051      //加密错误
	CITILEVEL_EXPIRE_TIME      = 86400     //城市分级数据失效时间
	LONLAT_EXPIRE_TIME         = 5 * 86400 //逆地理位置失效时间
	IPC_EXPIRE_TIME            = 86400     //高德ip数据缓存失效时间
	CURL_TIMEOUT               = 1         //生产环境curl超时
	CURL_TIMEOUT_TEST          = 2         //测试环境curl超时
	FLOATBIT                   = 2         //经纬度精度到小数点后多少位
	LOCAL_CACHE_EXPIRE_TIME    = 600
	AMAP_REST_API_GEOCODER_KEY = "99b7c33c53906ee93dc250e3d11bbb41" // "4c32ab0ed2b76d7f77c4c2d2523ac930" //
	HTTP_CHANNEL_BUFFER_LEN    = 10
)

var INFO = map[int]string{
	ERROR_PARAM:   "参数错误",
	ERROR_INTER:   "内部错误",
	ERROR_CAPTCHA: "验证码错误",
	ERROR_AUTH:    "鉴权失败",
}

const (
	PERPAGE    = 20
	JWT_SECRET = "XSpeUFjJ"
)
