package requests

import (
	"github.com/squzy/squzy_go/core"
	"net/http"
)

func NewRequest(trx *core.Transaction, req *http.Request) *http.Request {
	if trx == nil || trx.GetApplication() == nil {
		return req
	}
	// @TODO set meta here
	req.Header.Add(trx.GetApplication().GetTracingHeader(), trx.Id)
	return req
}