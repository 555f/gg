package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dto "github.com/f555/gg-examples/pkg/dto"
	errors "github.com/f555/gg-examples/pkg/errors"
	"io"
	"net/http"
)

type clientMethodOptions struct {
	ctx    context.Context
	before []ClientBeforeFunc
	after  []ClientAfterFunc
}
type ClientMethodOption func(*clientMethodOptions)

func WithContext(ctx context.Context) ClientMethodOption {
	return func(o *clientMethodOptions) {
		o.ctx = ctx
	}
}
func Before(before ...ClientBeforeFunc) ClientMethodOption {
	return func(o *clientMethodOptions) {
		o.before = append(o.before, before...)
	}
}
func After(after ...ClientAfterFunc) ClientMethodOption {
	return func(o *clientMethodOptions) {
		o.after = append(o.after, after...)
	}
}

type ClientBeforeFunc func(context.Context, *http.Request) context.Context
type ClientAfterFunc func(context.Context, *http.Response) context.Context
type ProfileControllerClient struct {
	client *http.Client
	target string
}
type ProfileControllerCreateRequest struct {
	c          *ProfileControllerClient
	client     *http.Client
	methodOpts *clientMethodOptions
	params     struct {
		firstName string
		lastName  string
		address   *string
		zip       int
	}
}

func (create *ProfileControllerCreateRequest) SetAddress(address string) *ProfileControllerCreateRequest {
	create.params.address = &address
	return create
}
func (create *ProfileControllerCreateRequest) Execute(opts ...ClientMethodOption) (profile *dto.Profile, err error) {
	for _, o := range opts {
		o(create.methodOpts)
	}
	var body struct {
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Address   *string `json:"address,omitempty"`
		Zip       int     `json:"zip"`
	}
	ctx, cancel := context.WithCancel(create.methodOpts.ctx)
	path := fmt.Sprintf("/profiles")
	req, err := http.NewRequest("POST", create.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	body.FirstName = create.params.firstName
	body.LastName = create.params.lastName
	body.Address = create.params.address
	body.Zip = create.params.zip
	var reqData bytes.Buffer
	err = json.NewEncoder(&reqData).Encode(body)
	if err != nil {
		cancel()
		return
	}
	req.Body = io.NopCloser(&reqData)

	req.Header.Add("Content-Type", "application/json")
	for _, before := range create.methodOpts.before {
		ctx = before(ctx, req)
	}
	resp, err := create.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	for _, after := range create.methodOpts.after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		err = json.NewDecoder(resp.Body).Decode(&errorWrapper)
		if err != nil {
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	var respBody struct {
		Profile *dto.Profile `json:"profile"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return
	}
	return respBody.Profile, err
}

type ProfileControllerDownloadFileRequest struct {
	c          *ProfileControllerClient
	client     *http.Client
	methodOpts *clientMethodOptions
	params     struct {
		id string
	}
}

func (downloadFile *ProfileControllerDownloadFileRequest) Execute(opts ...ClientMethodOption) (data string, err error) {
	for _, o := range opts {
		o(downloadFile.methodOpts)
	}
	ctx, cancel := context.WithCancel(downloadFile.methodOpts.ctx)
	path := fmt.Sprintf("/profiles/%s/file", downloadFile.params.id)
	req, err := http.NewRequest("GET", downloadFile.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	req.Header.Add("Content-Type", "application/json")
	for _, before := range downloadFile.methodOpts.before {
		ctx = before(ctx, req)
	}
	resp, err := downloadFile.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	for _, after := range downloadFile.methodOpts.after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		err = json.NewDecoder(resp.Body).Decode(&errorWrapper)
		if err != nil {
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	var respBody struct {
		Data string `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return
	}
	return respBody.Data, err
}

type ProfileControllerRemoveRequest struct {
	c          *ProfileControllerClient
	client     *http.Client
	methodOpts *clientMethodOptions
	params     struct {
		id string
	}
}

func (remove *ProfileControllerRemoveRequest) Execute(opts ...ClientMethodOption) (err error) {
	for _, o := range opts {
		o(remove.methodOpts)
	}
	ctx, cancel := context.WithCancel(remove.methodOpts.ctx)
	path := fmt.Sprintf("/profiles/%s", remove.params.id)
	req, err := http.NewRequest("DELETE", remove.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	req.Header.Add("Content-Type", "application/json")
	for _, before := range remove.methodOpts.before {
		ctx = before(ctx, req)
	}
	resp, err := remove.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	for _, after := range remove.methodOpts.after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		err = json.NewDecoder(resp.Body).Decode(&errorWrapper)
		if err != nil {
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	var respBody struct{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return
	}
	return err
}
func (profileController *ProfileControllerClient) Create(firstName string, lastName string, zip int) *ProfileControllerCreateRequest {
	m := &ProfileControllerCreateRequest{client: profileController.client, methodOpts: &clientMethodOptions{ctx: context.TODO()}, c: profileController}
	m.params.firstName = firstName
	m.params.lastName = lastName
	m.params.zip = zip
	return m
}
func (profileController *ProfileControllerClient) DownloadFile(id string) *ProfileControllerDownloadFileRequest {
	m := &ProfileControllerDownloadFileRequest{client: profileController.client, methodOpts: &clientMethodOptions{ctx: context.TODO()}, c: profileController}
	m.params.id = id
	return m
}
func (profileController *ProfileControllerClient) Remove(id string) *ProfileControllerRemoveRequest {
	m := &ProfileControllerRemoveRequest{client: profileController.client, methodOpts: &clientMethodOptions{ctx: context.TODO()}, c: profileController}
	m.params.id = id
	return m
}
func NewProfileControllerClient(target string) *ProfileControllerClient {
	c := &ProfileControllerClient{target: target, client: http.DefaultClient}
	return c
}
