package general_errors

import (
	"github.com/eden-framework/courier/status_error"
)

func init() {
	status_error.StatusErrorCodes.Register("NotFound", 404101000, "未找到", "", false)
	status_error.StatusErrorCodes.Register("Forbidden", 403101000, "不允许操作", "", true)
	status_error.StatusErrorCodes.Register("InternalError", 500101000, "内部处理错误", "", false)
	status_error.StatusErrorCodes.Register("BadRequest", 400101000, "请求参数错误", "", false)
	status_error.StatusErrorCodes.Register("Unauthorized", 401101000, "未授权", "", true)
	status_error.StatusErrorCodes.Register("Conflict", 409101000, "操作冲突", "", true)
	status_error.StatusErrorCodes.Register("BadGateway", 502101000, "上游错误", "", false)

}
