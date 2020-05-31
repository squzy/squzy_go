package requests

import (
	"context"
	"github.com/squzy/squzy_go/core"
	"net/http"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func NewRequest(trx *core.Transaction, req *http.Request) *http.Request {
	if trx == nil || trx.GetApplication() == nil {
		return req
	}
	trx.SetMeta()
	req.Header.Add(trx.GetApplication().GetTracingHeader(), trx.Id)
	return req.WithContext(context.WithValue(req.Context(), core.CONTEXT_KEY, trx))
}

func NewRoundTripper(parent http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		if nil == parent {
			parent = http.DefaultTransport
		}

		trx := core.GetTransactionFromContext(request.Context())
		request.Header.Add(trx.GetApplication().GetTracingHeader(), trx.Id)
		response, err := parent.RoundTrip(request)
		//@TODO set meta data
		trx.End()
		return response, err
	})
}