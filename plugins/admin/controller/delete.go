package controller

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/guard"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/response"
)

// Delete delete the row from database.
func (h *Handler) Delete(ctx *context.Context) {

	param := guard.GetDeleteParam(ctx)

	//token := ctx.FormValue("_t")
	//
	//if !auth.TokenHelper.CheckToken(token) {
	//	ctx.SetStatusCode(http.StatusBadRequest)
	//	ctx.WriteString(`{"code":400, "msg":"delete fail"}`)
	//	return
	//}

	if err := h.table(param.Prefix, ctx).DeleteData(param.Id); err != nil {
		logger.ErrorCtx(ctx, "Delete error %+v", err)
		response.Error(ctx, "delete fail")
		return
	}

	response.OkWithData(ctx, map[string]interface{}{
		"token": h.authSrv().AddToken(),
	})
}
