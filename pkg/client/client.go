package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dto "github.com/555f/gg/examples/rest-service/pkg/dto"
	errors "github.com/555f/gg/examples/rest-service/pkg/errors"
	"io"
	"net/http"
)

type ClientBeforeFunc func(context.Context, *http.Request) (context.Context, error)
type ClientAfterFunc func(context.Context, *http.Response) context.Context
type clientOptions struct {
	ctx    context.Context
	before []ClientBeforeFunc
	after  []ClientAfterFunc
}
type ClientOption func(*clientOptions)

func WithContext(ctx context.Context) ClientOption {
	return func(o *clientOptions) {
		o.ctx = ctx
	}
}
func Before(before ...ClientBeforeFunc) ClientOption {
	return func(o *clientOptions) {
		o.before = append(o.before, before...)
	}
}
func After(after ...ClientAfterFunc) ClientOption {
	return func(o *clientOptions) {
		o.after = append(o.after, after...)
	}
}

type ProfileControllerClient struct {
	client *http.Client
	target string
	opts   *clientOptions
}
type ProfileControllerCreateRequest struct {
	c      *ProfileControllerClient
	client *http.Client
	opts   *clientOptions
	params struct {
		firstName string
		lastName  string
		address   *string
		zip       int
	}
}

func (createRequest *ProfileControllerCreateRequest) SetAddress(address string) *ProfileControllerCreateRequest {
	createRequest.params.address = &address
	return createRequest
}
func (createRequest *ProfileControllerCreateRequest) Execute(opts ...ClientOption) (profile *dto.Profile, err error) {
	for _, o := range opts {
		o(createRequest.opts)
	}
	var body struct {
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Address   *string `json:"address,omitempty"`
		Zip       int     `json:"zip"`
	}
	ctx, cancel := context.WithCancel(createRequest.opts.ctx)
	path := "/profiles"
	req, err := http.NewRequest("POST", createRequest.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	body.FirstName = createRequest.params.firstName
	body.LastName = createRequest.params.lastName
	body.Address = createRequest.params.address
	body.Zip = createRequest.params.zip
	var reqData bytes.Buffer
	err = json.NewEncoder(&reqData).Encode(body)
	if err != nil {
		cancel()
		return
	}
	req.Body = io.NopCloser(&reqData)

	req.Header.Add("Content-Type", "application/json")
	before := append(createRequest.c.opts.before, createRequest.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := createRequest.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(createRequest.c.opts.after, createRequest.opts.after...)
	for _, after := range after {
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
	return respBody.Profile, nil
}

type ProfileControllerDownloadFileRequest struct {
	c      *ProfileControllerClient
	client *http.Client
	opts   *clientOptions
	params struct {
		id string
	}
}

func (downloadFileRequest *ProfileControllerDownloadFileRequest) Execute(opts ...ClientOption) (data string, err error) {
	for _, o := range opts {
		o(downloadFileRequest.opts)
	}
	ctx, cancel := context.WithCancel(downloadFileRequest.opts.ctx)
	path := fmt.Sprintf("/profiles/%s/file", downloadFileRequest.params.id)
	req, err := http.NewRequest("GET", downloadFileRequest.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	req.Header.Add("Content-Type", "application/json")
	before := append(downloadFileRequest.c.opts.before, downloadFileRequest.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := downloadFileRequest.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(downloadFileRequest.c.opts.after, downloadFileRequest.opts.after...)
	for _, after := range after {
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
	return respBody.Data, nil
}

type ProfileControllerRemoveRequest struct {
	c      *ProfileControllerClient
	client *http.Client
	opts   *clientOptions
	params struct {
		id string
	}
}

func (removeRequest *ProfileControllerRemoveRequest) Execute(opts ...ClientOption) (err error) {
	for _, o := range opts {
		o(removeRequest.opts)
	}
	var body struct {
		Id string `json:"id"`
	}
	ctx, cancel := context.WithCancel(removeRequest.opts.ctx)
	path := "/profiles/{id}"
	req, err := http.NewRequest("DELETE", removeRequest.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	body.Id = removeRequest.params.id
	var reqData bytes.Buffer
	err = json.NewEncoder(&reqData).Encode(body)
	if err != nil {
		cancel()
		return
	}
	req.Body = io.NopCloser(&reqData)

	req.Header.Add("Content-Type", "application/json")
	before := append(removeRequest.c.opts.before, removeRequest.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := removeRequest.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(removeRequest.c.opts.after, removeRequest.opts.after...)
	for _, after := range after {
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
	return nil
}
func (profileController *ProfileControllerClient) Create(firstName string, lastName string, zip int) *ProfileControllerCreateRequest {
	m := &ProfileControllerCreateRequest{client: profileController.client, opts: &clientOptions{ctx: context.TODO()}, c: profileController}
	m.params.firstName = firstName
	m.params.lastName = lastName
	m.params.zip = zip
	return m
}
func (profileController *ProfileControllerClient) DownloadFile(id string) *ProfileControllerDownloadFileRequest {
	m := &ProfileControllerDownloadFileRequest{client: profileController.client, opts: &clientOptions{ctx: context.TODO()}, c: profileController}
	m.params.id = id
	return m
}
func (profileController *ProfileControllerClient) Remove(id string) *ProfileControllerRemoveRequest {
	m := &ProfileControllerRemoveRequest{client: profileController.client, opts: &clientOptions{ctx: context.TODO()}, c: profileController}
	m.params.id = id
	return m
}
func NewProfileControllerClient(target string, opts ...ClientOption) *ProfileControllerClient {
	c := &ProfileControllerClient{target: target, client: http.DefaultClient, opts: &clientOptions{}}
	for _, o := range opts {
		o(c.opts)
	}
	return c
}
