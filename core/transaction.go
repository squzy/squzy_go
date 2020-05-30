package core

import (
	gonanoid "github.com/matoous/go-nanoid"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Transactor interface {
	CreateTransaction(name string, trType api.TransactionType, parent *Transaction) (*Transaction, error)
}

type Transaction struct {
	Id string
	Name string
	Type api.TransactionType
	Parent *Transaction
	app *Application
}

func (t *Transaction) SetMeta() {

}

func (t *Transaction) End() {

}

func (t *Transaction) GetApplication() *Application {
	return t.app
}

func (t *Transaction) CreateTransaction(name string, trType api.TransactionType, parent *Transaction) (*Transaction, error) {
	return createTransaction(name, trType, t, t.app)
}

func createTransaction(name string, trType api.TransactionType, parent *Transaction, application *Application) (*Transaction, error) {
	id , err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		Id:   id,
		Name: name,
		Type: trType,
		Parent: parent,
		app: application,
	}, nil
}