// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package page

import (
	"bytes"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// SetPageContent set and return the panel of page content.
func SetPageContent(ctx *context.Context, user models.UserModel, c func(ctx interface{}) (types.Panel, error), conn db.Connection) {

	panel, err := c(ctx)

	if err != nil {
		logger.ErrorCtx(ctx, "SetPageContent %+v", err)
		panel = template.WarningPanel(ctx, err.Error())
	}

	tmpl, tmplName := template.Get(ctx, config.GetTheme()).GetTemplate(ctx.IsPjax())

	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, conn, ctx.Lang()).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Panel:        panel.GetContent(config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(ctx),
		TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
		TmplFootJS:   template.Default(ctx).GetFootJS(),
		Iframe:       ctx.IsIframe(),
	}))
	if err != nil {
		logger.ErrorCtx(ctx, "SetPageContent %+v", err)
	}
	ctx.WriteString(buf.String())
}
