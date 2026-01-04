package controller

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/guard"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/response"
)

// Update update the table row of given id.
func (h *Handler) Update(ctx *context.Context) {

	param := guard.GetUpdateParam(ctx)

	err := param.Panel.UpdateData(ctx, param.Value)

	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Ok(ctx)
}
