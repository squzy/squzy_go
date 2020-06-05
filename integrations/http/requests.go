package squzy_http

import (
	"context"
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

func NewRoundTripper(parent http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		if nil == parent {
			parent = http.DefaultTransport
		}

		trx := core.GetTransactionFromContext(request.Context())
		request.Header.Add(trx.GetApplication().GetTracingHeader(), trx.GetId())
		response, err := parent.RoundTrip(request)
		var path string

		if request.URL != nil {
			path = request.URL.Path
		}
		trx.SetMeta(&core.TransactionMeta{
			Host:   request.Host,
			Path:   path,
			Method: request.Method,
		}).End(err)

		return response, err
	})
}
