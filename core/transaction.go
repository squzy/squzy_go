package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"net/http"
	"time"
)

const (
	CONTEXT_KEY = "__squzy_transaction"
)

type Transactor interface {
	CreateTransaction(name string, trType api.TransactionType, parent *Transaction) *Transaction
}

type transactionRequestMsg struct {
	Id       string                `json:"id"`
	ParentId string                `json:"string,omitempty"`
	Name     string                `json:"name"`
	DateFrom string                `json:"dateFrom"`
	DateTo   string                `json:"dateTo"`
	Status   api.TransactionStatus `json:"status"`
	Type     api.TransactionType   `json:"type"`
	Error    *trxError             `json:"error,omitempty"`
	Meta     *TransactionMeta      `json:"meta,omitempty"`
}

type trxError struct {
	Message string `json:"message"`
}

type TransactionMeta struct {
	Host   string `json:"host"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type Transaction struct {
	Id         string
	Name       string
	Type       api.TransactionType
	Parent     *Transaction
	app        *Application
	httpClient *http.Client
	startTime  time.Time
	endTime    time.Time
	meta       *TransactionMeta
}

func (t *Transaction) SetMeta(meta *TransactionMeta) *Transaction {
	if t == nil {
		return nil
	}
	t.meta = meta
	return t
}

func (t *Transaction) End(err error) {
	if t == nil {
		return
	}
	t.endTime = time.Now()
	var status = api.TransactionStatus_TRANSACTION_SUCCESSFUL
	var trErr *trxError
	if err != nil {
		status = api.TransactionStatus_TRANSACTION_FAILED
		trErr = &trxError{
			Message: err.Error(),
		}
	}

	rq := &transactionRequestMsg{
		Id:       t.Id,
		ParentId: t.getParentId(),
		Name:     t.Name,
		DateFrom: t.startTime.Format(time.RFC3339),
		DateTo:   t.endTime.Format(time.RFC3339),
		Status:   status,
		Type:     t.Type,
		Error:    trErr,
		Meta:     t.meta,
	}

	b, err := json.Marshal(rq)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/applications/%s/transactions", t.GetApplication().monitoringHost, t.GetApplication().GetID()), bytes.NewReader(b))
	if err != nil {
		return
	}
	_, _ = sendHttp(t.GetApplication().GetHttpClient(), req)
}

func (t *Transaction) GetApplication() *Application {
	if t == nil {
		return nil
	}
	return t.app
}

func (t *Transaction) CreateTransaction(name string, trType api.TransactionType, parent *Transaction) *Transaction {
	if t == nil {
		return nil
	}
	return createTransaction(name, trType, t, t.app)
}

func (t *Transaction) getParentId() string {
	if t == nil {
		return ""
	}
	if t.Parent != nil {
		return t.Parent.Id
	}
	return ""
}

func createTransaction(name string, trType api.TransactionType, parent *Transaction, application *Application) *Transaction {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil
	}
	return &Transaction{
		Id:        id,
		Name:      name,
		Type:      trType,
		Parent:    parent,
		app:       application,
		startTime: time.Now(),
	}
}

func GetTransactionFromContext(ctx context.Context) *Transaction {
	trx, ok := ctx.Value(CONTEXT_KEY).(*Transaction)
	if ok {
		return trx
	}
	return nil
}

func ContextWithTransaction(ctx context.Context, trx *Transaction) context.Context {
	return context.WithValue(ctx, CONTEXT_KEY, trx)
}
