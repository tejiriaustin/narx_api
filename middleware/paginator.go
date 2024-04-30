package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tejiriaustin/narx_api/constants"
	"github.com/tejiriaustin/narx_api/response"
)

func ReadPaginationOptions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := ctx.Params

		if page, _ := params.Get("page"); page != "" {
			pageNum, err := strconv.ParseInt(page, 10, 64)
			if err != nil {
				response.FormatResponse(ctx, http.StatusBadRequest, "page number must be a number", nil)
				return
			}

			ctx.Set(string(constants.ContextKeyPageNumber), pageNum)
		}

		if perPage, _ := params.Get("per_page"); perPage != "" {
			perPageNum, err := strconv.ParseInt(perPage, 10, 64)
			if err != nil {
				response.FormatResponse(ctx, http.StatusBadRequest, "per page number must be a number", nil)
				return
			}

			ctx.Set(string(constants.ContextKeyPageNumber), perPageNum)
		}

	}
}
