package squzy_http

import (
	"context"
	"fmt"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
	"net/http"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return fn(r) }

func NewRequest(trx *core.Transaction, req *http.Request) *http.Request {
	if trx == nil || trx.GetApplication() == nil {
		return req
	}
	return req.WithContext(context.WithValue(req.Context(), core.CONTEXT_KEY, trx))
}

func NewRoundTripper(app *core.Application, parent http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		if nil == parent {
			parent = http.DefaultTransport
		}
		var path string

		if request.URL != nil {
			path = request.URL.Path
		}

		trx := core.GetTransactionFromContext(request.Context())
		if trx == nil {
			trx = app.CreateTransaction(fmt.Sprintf("%s/%s", request.Host, path), api.TransactionType_TRANSACTION_TYPE_HTTP, nil)
		}
		request.Header.Add(app.GetTracingHeader(), trx.GetId())
		response, err := parent.RoundTrip(request)
		trx.SetMeta(&core.TransactionMeta{
			Host:   request.Host,
			Path:   path,
			Method: request.Method,
		}).End(err)

		return response, err
	})
}
