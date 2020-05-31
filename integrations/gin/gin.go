package squyz_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/squzy/squzy_go/core"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
)

func New(app *core.Application) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()

		trx := app.CreateTransaction(path, api.TransactionType_TRANSACTION_TYPE_ROUTER, &core.Transaction{
			Id: app.GetID(),
		})

		context.Set(core.CONTEXT_KEY, trx)
		// @TODO set meta
		context.Next()

		trx.End()
	}
}