// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/miniflux/miniflux2/model"
	"github.com/miniflux/miniflux2/server/core"
)

// ShowHistoryPage renders the page with all read entries.
func (c *Controller) ShowHistoryPage(ctx *core.Context, request *core.Request, response *core.Response) {
	user := ctx.LoggedUser()
	offset := request.QueryIntegerParam("offset", 0)

	args, err := c.getCommonTemplateArgs(ctx)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	builder := c.store.GetEntryQueryBuilder(user.ID, user.Timezone)
	builder.WithStatus(model.EntryStatusRead)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(model.DefaultSortingDirection)
	builder.WithOffset(offset)
	builder.WithLimit(NbItemsPerPage)

	entries, err := builder.GetEntries()
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	count, err := builder.CountEntries()
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.HTML().Render("history", args.Merge(tplParams{
		"entries":    entries,
		"total":      count,
		"pagination": c.getPagination(ctx.Route("history"), count, offset),
		"menu":       "history",
	}))
}

// FlushHistory changes all "read" items to "removed".
func (c *Controller) FlushHistory(ctx *core.Context, request *core.Request, response *core.Response) {
	user := ctx.LoggedUser()

	err := c.store.FlushHistory(user.ID)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.Redirect(ctx.Route("history"))
}
