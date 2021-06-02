package public

const (
	ValidatorKey         = "ValidatorKey"
	TranslatorKey        = "TranslatorKey"
	AdminLoginSessionKey = "AdminLoginSessionInfo"

	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	HTTPRuleNeedHttps = 1

	HTTPAccessControlAuth = 1

	IsDeleteTrue  = 1
	IsDeleteFalse = 0

	RedisFlowDayKey = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal          = "flow_total"
	FlowServicePrefix  = "flow_service_"
	FlowAppPrefix = "flow_app_"

	JwtSignKey = "a00b66bbe6cb4961a9d87503c8556bab"
	JwtExpires = 60 * 60
	TokenType = "Bearer"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)