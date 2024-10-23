package route

var (
	ROUTE_USER_REGISTER = "/user/register"
	ROUTE_USER_LOGIN    = "/user/login"
	ROUTE_USER_DELETE   = "/user/{id}"
)

var ROUTER_NEED_ADMIN_AUTH = []string{ROUTE_USER_DELETE}
var ROUTER_NEED_AUTH = append([]string{}, ROUTER_NEED_ADMIN_AUTH...)
