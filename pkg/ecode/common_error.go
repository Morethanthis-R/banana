package ecode

// All dictionary ecode
var (
	OK           = add(0)        //请求成功
	INTERNAL_ERR = add(10000000) //系统内部错误
	SYS_PAUSE    = add(10000001) //系统暂停服务
	IP_LIMIT     = add(10000002) //IP访问受限
	SYS_BUSY     = add(10000003) //系统繁忙

	REQUEST_EXCEPTION = add(10001000) //请求异常
	LACK_OF_PARAM     = add(10001001) //缺乏必要参数
	INVALID_PARAM     = add(10001002) //参数不符合规则

	RATE_LIMIT         = add(10001003) //请求频率超出限制
	UNSUPPORTED_FORMAT = add(10001004) //不支持的数据格式
	SIZE_LIMIT         = add(10001005) //请求超出大小限制
	NOT_FOUND          = add(10001006) //访问资源不存在x
	INVALID_DATE_PARAM = add(10001007) //参数日期区间有误

	AUTH_TOKEN_NOT_FOUND = add(10002000) //鉴权标志不存在
	AUTH_TOKEN_EXPIRED   = add(10002001) //鉴权标志已过期
	AUTH_TOKEN_INVALID   = add(10002002) //无效的鉴权标志
	PERMISSION_DENIED    = add(10002003) //无权限 访问被拒绝
	AUTH_FAIL            = add(10002004) //传入参数不对导致授权失败

	MYSQL_CONN_ERR     = add(10003000) //MySQL连接失败
	MYSQL_CONN_EXPIRED = add(10003001) //MySQL连接失效
	MYSQL_ERR          = add(10003002) //MySQL操作异常

	REDIS_CONN_ERR     = add(10004000) //Redis连接失败
	REDIS_CONN_EXIPRED = add(10004001) //Redis连接失效
	REDIS_ERR          = add(10004002) //Redis操作异常

	ES_CONN_ERR     = add(10005000) //ES连接错误
	ES_CONN_EXPIRED = add(10005001) //ES连接失败
	ES_ERR          = add(10005002) //ES操作异常
	ES_INDEX_ERROR  = add(10005003) //ES索引异常
	ES_QUERY_ERROR  = add(10005004) //ES查询异常

	EXTERNAL_API_PARAM_ERR   = add(10010001) //访问外部api参数错误
	EXTERNAL_API_FAIL        = add(10010002) //访问外部api失败
	EXTERNAL_API_NO_RESPONSE = add(10010003) //访问外部api无响应
	EXTERNAL_API_JSON_ERR    = add(10010004) //访问外部api返回内容json解析错误

	FILE_NOT_EXIST    = add(10011001) //文件不存在
	FILE_TYPE_ERR     = add(10011002) //文件类型/格式错误
	FILE_FORMAT_ERR   = add(10011003) //文件内部格式错误
	FILE_OVERSIZE     = add(10011004) //文件大小超出限制
	FILE_PARSE_ERR    = add(10011005) //文件解析错误
	FILE_TRANSFER_ERR = add(10011006) //文件传输错误

)

var CommonErrMsg = map[int]string{

	OK.Code():           "请求成功",
	INTERNAL_ERR.Code(): "系统内部错误",
	SYS_PAUSE.Code():    "系统暂停服务",
	IP_LIMIT.Code():     "IP访问受限",
	SYS_BUSY.Code():     "系统繁忙",

	REQUEST_EXCEPTION.Code():  "请求异常",
	LACK_OF_PARAM.Code():      "缺乏必要参数",
	INVALID_PARAM.Code():      "参数不符合规则",
	INVALID_DATE_PARAM.Code(): "参数日期区间有误",
	RATE_LIMIT.Code():         "请求频率超出限制",
	UNSUPPORTED_FORMAT.Code(): "不支持的数据格式",
	SIZE_LIMIT.Code():         "请求超出大小限制",
	NOT_FOUND.Code():          "访问资源不存在",

	AUTH_TOKEN_NOT_FOUND.Code(): "鉴权标志不存在",
	AUTH_TOKEN_EXPIRED.Code():   "鉴权标志已过期",
	AUTH_TOKEN_INVALID.Code():   "无效的鉴权标志",
	PERMISSION_DENIED.Code():    "无权限 访问被拒绝",
	AUTH_FAIL.Code():            "授权失败",

	MYSQL_CONN_ERR.Code():     "MySQL连接失败",
	MYSQL_CONN_EXPIRED.Code(): "MySQL连接失效",
	MYSQL_ERR.Code():          "MySQL操作异常",

	REDIS_CONN_ERR.Code():     "Redis连接失败",
	REDIS_CONN_EXIPRED.Code(): "Redis连接失效",
	REDIS_ERR.Code():          "Redis操作异常",

	ES_CONN_ERR.Code():     "ES连接错误",
	ES_CONN_EXPIRED.Code(): "ES连接失效",
	ES_ERR.Code():          "ES操作异常",
	ES_INDEX_ERROR.Code():  "ES索引异常",
	ES_QUERY_ERROR.Code():  "ES查询异常",

	EXTERNAL_API_PARAM_ERR.Code():   "访问外部api参数错误",
	EXTERNAL_API_FAIL.Code():        "访问外部api失败",
	EXTERNAL_API_NO_RESPONSE.Code(): "访问外部api无响应",
	EXTERNAL_API_JSON_ERR.Code():    "访问外部api返回内容json解析错误",

	FILE_NOT_EXIST.Code():    "文件不存在",
	FILE_TYPE_ERR.Code():     "文件类型/格式错误",
	FILE_FORMAT_ERR.Code():   "文件内部格式错误",
	FILE_OVERSIZE.Code():     "文件大小超出限制",
	FILE_PARSE_ERR.Code():    "文件解析错误",
	FILE_TRANSFER_ERR.Code(): "文件传输错误",
}
