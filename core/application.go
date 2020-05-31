package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"io/ioutil"
	"net/http"
)

type Application struct {
	id string
	monitoringHost string
	tracingHeader string
}

type Options struct {
	MonitoringHost string
	ApplicationName string
	ApplicationHost string
}

func (a *Application) CreateTransaction(name string, trType api.TransactionType, parent *Transaction) (*Transaction) {
	if a == nil {
		return nil
	}
	return createTransaction(name, trType, nil, a)
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

func CreateApplication(client *http.Client, opts *Options) (*Application, error) {
	type reqBody struct {
		name string `json:"name"`
		host string `json:"host"`
	}
	req := &reqBody{
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

	response, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	if response != nil {
		defer response.Body.Close()
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	type res struct {
		Data struct{
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
		id: responseJson.Data.ApplicationID,
		monitoringHost: opts.MonitoringHost,
		tracingHeader: responseJson.Data.TracingHeader,
	}, nil
}