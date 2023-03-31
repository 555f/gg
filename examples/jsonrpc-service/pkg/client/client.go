package client

import (
	"bytes"
	"context"
	"encoding/json"
	dto "github.com/f555/gg-examples/pkg/dto"
	"io"
	"net/http"
	"sync/atomic"
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
type ClientAfterFunc func(context.Context, *http.Response, json.RawMessage) context.Context
type clientRequester interface {
	makeRequest() (string, any)
	makeResult(data []byte) (any, error)
	before() []ClientBeforeFunc
	after() []ClientAfterFunc
	context() context.Context
}
type BatchResult struct {
	results []any
}

func (r *BatchResult) At(i int) any {
	return r.results[i]
}
func (r *BatchResult) Len() int {
	return len(r.results)
}

type clientReq struct {
	ID      uint64 `json:"id"`
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}
type clientResp struct {
	ID      uint64          `json:"id"`
	Version string          `json:"jsonrpc"`
	Error   *clientError    `json:"error"`
	Result  json.RawMessage `json:"result"`
}
type clientError struct {
	Code    error  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type ProfileControllerClient struct {
	client      *http.Client
	target      string
	incrementID uint64
	methodOpts  *clientMethodOptions
}

func (profileController *ProfileControllerClient) autoIncrementID() uint64 {
	return atomic.AddUint64(&profileController.incrementID, 1)
}
func (profileController *ProfileControllerClient) Execute(requests ...clientRequester) (*BatchResult, error) {
	profileController.incrementID = 0
	req, err := http.NewRequest("POST", profileController.target, nil)
	if err != nil {
		return nil, err
	}
	idsIndex := make(map[uint64]int, len(requests))
	rpcRequests := make([]clientReq, len(requests))
	for _, beforeFunc := range profileController.methodOpts.before {
		req = req.WithContext(beforeFunc(req.Context(), req))
	}
	for i, request := range requests {
		req = req.WithContext(request.context())
		for _, beforeFunc := range request.before() {
			req = req.WithContext(beforeFunc(req.Context(), req))
		}
		methodName, params := request.makeRequest()
		r := clientReq{ID: profileController.autoIncrementID(), Version: "2.0", Method: methodName, Params: params}
		idsIndex[r.ID] = i
		rpcRequests[i] = r
	}
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(rpcRequests); err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(buf)
	resp, err := profileController.client.Do(req)
	if err != nil {
		return nil, err
	}
	responses := make([]clientResp, len(requests))
	if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
		return nil, err
	}
	batchResult := &BatchResult{results: make([]any, len(requests))}
	for _, response := range responses {
		for _, afterFunc := range profileController.methodOpts.after {
			afterFunc(resp.Request.Context(), resp, response.Result)
		}
		i := idsIndex[response.ID]
		request := requests[i]
		for _, afterFunc := range request.after() {
			afterFunc(resp.Request.Context(), resp, response.Result)
		}
		result, err := request.makeResult(response.Result)
		if err != nil {
			return nil, err
		}
		batchResult.results[i] = result
	}
	return batchResult, nil
}

type ProfileControllerCreateBatchResult struct {
	Profile *dto.Profile `json:"profile"`
}
type ProfileControllerCreateRequest struct {
	c          *ProfileControllerClient
	client     *http.Client
	methodOpts *clientMethodOptions
	params     struct {
		firstName string
		lastName  string
		address   *string
	}
}

func (create *ProfileControllerCreateRequest) SetAddress(address string) *ProfileControllerCreateRequest {
	create.params.address = &address
	return create
}
func (create *ProfileControllerCreateRequest) makeRequest() (string, any) {
	var params struct {
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Address   *string `json:"address,omitempty"`
	}
	params.FirstName = create.params.firstName
	params.LastName = create.params.lastName
	params.Address = create.params.address
	return "profile.create", params
}
func (create *ProfileControllerCreateRequest) makeResult(data []byte) (any, error) {
	var result ProfileControllerCreateBatchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
func (create *ProfileControllerCreateRequest) before() []ClientBeforeFunc {
	return create.methodOpts.before
}
func (create *ProfileControllerCreateRequest) after() []ClientAfterFunc {
	return create.methodOpts.after
}
func (create *ProfileControllerCreateRequest) context() context.Context {
	return create.methodOpts.ctx
}
func (create *ProfileControllerCreateRequest) Execute() (profile *dto.Profile, err error) {
	batchResult, err := create.c.Execute(create)
	if err != nil {
		return
	}
	clientResult := batchResult.At(0).(ProfileControllerCreateBatchResult)
	return clientResult.Profile, err
}

type ProfileControllerRemoveRequest struct {
	c          *ProfileControllerClient
	client     *http.Client
	methodOpts *clientMethodOptions
	params     struct {
		id string
	}
}

func (remove *ProfileControllerRemoveRequest) makeRequest() (string, any) {
	var params struct {
		Id string `json:"id"`
	}
	params.Id = remove.params.id
	return "profile.delete", params
}
func (remove *ProfileControllerRemoveRequest) makeResult(data []byte) (any, error) {
	return nil, nil
}
func (remove *ProfileControllerRemoveRequest) before() []ClientBeforeFunc {
	return remove.methodOpts.before
}
func (remove *ProfileControllerRemoveRequest) after() []ClientAfterFunc {
	return remove.methodOpts.after
}
func (remove *ProfileControllerRemoveRequest) context() context.Context {
	return remove.methodOpts.ctx
}
func (remove *ProfileControllerRemoveRequest) Execute() (err error) {
	_, err = remove.c.Execute(remove)
	if err != nil {
		return
	}
	return err
}
func (profileController *ProfileControllerClient) Create(firstName string, lastName string, opts ...ClientMethodOption) *ProfileControllerCreateRequest {
	m := &ProfileControllerCreateRequest{client: profileController.client, methodOpts: &clientMethodOptions{ctx: context.TODO()}, c: profileController}
	m.params.firstName = firstName
	m.params.lastName = lastName
	for _, o := range opts {
		o(m.methodOpts)
	}
	return m
}
func (profileController *ProfileControllerClient) Remove(id string, opts ...ClientMethodOption) *ProfileControllerRemoveRequest {
	m := &ProfileControllerRemoveRequest{client: profileController.client, methodOpts: &clientMethodOptions{ctx: context.TODO()}, c: profileController}
	m.params.id = id
	for _, o := range opts {
		o(m.methodOpts)
	}
	return m
}
func NewProfileControllerClient(target string, opts ...ClientMethodOption) *ProfileControllerClient {
	c := &ProfileControllerClient{target: target, client: http.DefaultClient, methodOpts: &clientMethodOptions{ctx: context.TODO()}}
	for _, o := range opts {
		o(c.methodOpts)
	}
	return c
}
