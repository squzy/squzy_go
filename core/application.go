package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"net/http"
)

type Application struct {
	id            string
	apiHost       string
	host          string
	tracingHeader string
	httpClient    *http.Client
}

type Options struct {
	ApiHost         string
	ApplicationName string
	ApplicationHost string
}

func (a *Application) CreateTransaction(name string, trType api.TransactionType, parent *Transaction) *Transaction {
	if a == nil {
		return nil
	}
	return New(name, trType, a, parent)
}

func (a *Application) GetTracingHeader() string {
	if a == nil {
		return ""
	}
	return a.tracingHeader
}

func (a *Application) GetID() string {
	if a == nil {
		return ""
	}
	return a.id
}

func (a *Application) GetHost() string {
	if a == nil {
		return ""
	}
	return a.host
}

func (a *Application) GetHttpClient() *http.Client {
	if a == nil {
		return nil
	}
	return a.httpClient
}

func (a *Application) GetApiHost() string {
	if a == nil {
		return ""
	}
	return a.apiHost
}

type registerAppRequestBody struct {
	name string `json:"name"`
	host string `json:"host"`
}

func CreateApplication(client *http.Client, opts *Options) (*Application, error) {
	req := &registerAppRequestBody{
		name: opts.ApplicationName,
		host: opts.ApplicationHost,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/applications", opts.ApplicationHost), bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	bodyBytes, err := sendHttp(client, r)

	if err != nil {
		return nil, err
	}

	type res struct {
		Data struct {
			ApplicationID string `json:"application_id"`
			TracingHeader string `json:"tracing_header"`
		} `json:"data"`
	}

	responseJson := &res{}

	err = json.Unmarshal(bodyBytes, responseJson)

	if err != nil {
		return nil, err
	}
	return &Application{
		id:            responseJson.Data.ApplicationID,
		apiHost:       opts.ApiHost,
		host:          opts.ApplicationHost,
		tracingHeader: responseJson.Data.TracingHeader,
		httpClient:    client,
	}, nil
}
