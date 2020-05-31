package squyz_gin

import (
	"github.com/gin-gonic/gin"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
)

func New(app *core.Application) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()

		trx := app.CreateTransaction(path, api.TransactionType_TRANSACTION_TYPE_ROUTER, &core.Transaction{
			Id: context.GetHeader(app.GetTracingHeader()),
		})

		context.Set(core.CONTEXT_KEY, trx)

		context.Next()

		// @TODO set meta
		trx.SetMeta().End()
	}
}