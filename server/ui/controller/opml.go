// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package controller

import (
	"log"

	"github.com/miniflux/miniflux2/server/core"
)

func (c *Controller) Export(ctx *core.Context, request *core.Request, response *core.Response) {
	user := ctx.LoggedUser()
	opml, err := c.opmlHandler.Export(user.ID)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.XML().Download("feeds.opml", opml)
}

func (c *Controller) Import(ctx *core.Context, request *core.Request, response *core.Response) {
	args, err := c.getCommonTemplateArgs(ctx)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.HTML().Render("import", args.Merge(tplParams{
		"menu": "feeds",
	}))
}

func (c *Controller) UploadOPML(ctx *core.Context, request *core.Request, response *core.Response) {
	file, fileHeader, err := request.File("file")
	if err != nil {
		log.Println(err)
		response.Redirect(ctx.Route("import"))
		return
	}
	defer file.Close()

	user := ctx.LoggedUser()
	log.Printf("[UI:UploadOPML] User #%d uploaded this file: %s (%d bytes)\n", user.ID, fileHeader.Filename, fileHeader.Size)

	if impErr := c.opmlHandler.Import(user.ID, file); impErr != nil {
		args, err := c.getCommonTemplateArgs(ctx)
		if err != nil {
			response.HTML().ServerError(err)
			return
		}

		response.HTML().Render("import", args.Merge(tplParams{
			"errorMessage": impErr,
			"menu":         "feeds",
		}))

		return
	}

	response.Redirect(ctx.Route("feeds"))
}
