package core

import (
	"context"
	"errors"
	gonanoid "github.com/matoous/go-nanoid"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

const (
	CONTEXT_KEY = "__squzy_transaction"
)

var (
	errNotFound = errors.New("can't find transaction")
)

type Transactor interface {
	CreateTransaction(name string, trType api.TransactionType, parent *Transaction) *Transaction
}

type Transaction struct {
	Id        string
	Name      string
	Type      api.TransactionType
	Parent    *Transaction
	app       *Application
	startTime time.Time
	endTime   time.Time
}

func (t *Transaction) SetMeta() *Transaction {
	if t == nil {
		return nil
	}
	return t
}

func (t *Transaction) End() {
	if t == nil {
		return
	}
	t.endTime = time.Now()
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

func GetTransactionFromContext(ctx context.Context) (*Transaction) {
	trx, ok := ctx.Value(CONTEXT_KEY).(*Transaction)
	if ok {
		return trx
	}
	return nil
}

func ContextWithTransaction(ctx context.Context, trx *Transaction) context.Context {
	return context.WithValue(ctx, CONTEXT_KEY, trx)
}
