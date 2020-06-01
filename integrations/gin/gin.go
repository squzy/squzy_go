package squyz_gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
	"strings"
)

func New(app *core.Application) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()

		trx := app.CreateTransaction(path, api.TransactionType_TRANSACTION_TYPE_ROUTER, &core.Transaction{
			Id: context.GetHeader(app.GetTracingHeader()),
		})

		var method string
		if context.Request != nil {
			method = context.Request.Method
		}

		trx.SetMeta(&core.TransactionMeta{
			Path: path,
			Method: method,
		})

		context.Set(core.CONTEXT_KEY, trx)

		context.Next()
		var err error = nil

		if len(context.Errors) > 0 {
			err = errors.New(strings.Join(context.Errors.Errors(), ";"))
		}

		trx.End(err)
	}
}