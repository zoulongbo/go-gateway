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

	IsDeleteTrue  = 1
	IsDeleteFalse = 0

	FlowTotal          = "flow_total"
	FlowServicePrefix  = "flow_service_"
	FlowAppPrefix = "flow_app_"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)