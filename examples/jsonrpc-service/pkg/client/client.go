package client

import (
	"context"
	"encoding/json"
	dto "github.com/555f/gg/examples/jsonrpc-service/pkg/dto"
	jsonrpc "github.com/555f/jsonrpc"
	"net/http"
)

type ProfileControllerClient struct {
	*jsonrpc.Client
}
type ProfileControllerCreateBatchResult struct {
	Profile *dto.Profile `json:"profile"`
}
type ProfileControllerCreateRequest struct {
	c      *ProfileControllerClient
	params struct {
		token     string
		firstName string
		lastName  string
		address   *string
	}
	before []jsonrpc.ClientBeforeFunc
	after  []jsonrpc.ClientAfterFunc
	ctx    context.Context
}

func (create *ProfileControllerCreateRequest) SetAddress(address string) *ProfileControllerCreateRequest {
	create.params.address = &address
	return create
}
func (create *ProfileControllerCreateRequest) MakeRequest() (string, any) {
	var params struct {
		Token     string  `json:"token"`
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Address   *string `json:"address,omitempty"`
	}
	params.Token = create.params.token
	params.FirstName = create.params.firstName
	params.LastName = create.params.lastName
	params.Address = create.params.address
	return "profile.create", params
}
func (create *ProfileControllerCreateRequest) SetBefore(before ...jsonrpc.ClientBeforeFunc) *ProfileControllerCreateRequest {
	create.before = before
	return create
}
func (create *ProfileControllerCreateRequest) SetAfter(after ...jsonrpc.ClientAfterFunc) *ProfileControllerCreateRequest {
	create.after = after
	return create
}
func (create *ProfileControllerCreateRequest) WithContext(ctx context.Context) *ProfileControllerCreateRequest {
	create.ctx = ctx
	return create
}
func (create *ProfileControllerCreateRequest) RawExecute() ([]byte, map[uint64]int, *http.Response, error) {
	return create.c.RawExecute(create)
}
func (create *ProfileControllerCreateRequest) MakeResult(data []byte) (any, error) {
	var result ProfileControllerCreateBatchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
func (create *ProfileControllerCreateRequest) Before() []jsonrpc.ClientBeforeFunc {
	return create.before
}
func (create *ProfileControllerCreateRequest) After() []jsonrpc.ClientAfterFunc {
	return create.after
}
func (create *ProfileControllerCreateRequest) Context() context.Context {
	return create.ctx
}
func (create *ProfileControllerCreateRequest) Execute() (profile *dto.Profile, err error) {
	batchResult, err := create.c.Client.Execute(create)
	if err != nil {
		return
	}
	clientResult := batchResult.At(0).(ProfileControllerCreateBatchResult)
	return clientResult.Profile, err
}

type ProfileControllerRemoveRequest struct {
	c      *ProfileControllerClient
	params struct {
		id string
	}
	before []jsonrpc.ClientBeforeFunc
	after  []jsonrpc.ClientAfterFunc
	ctx    context.Context
}

func (remove *ProfileControllerRemoveRequest) MakeRequest() (string, any) {
	var params struct {
		Id string `json:"id"`
	}
	params.Id = remove.params.id
	return "profile.delete", params
}
func (remove *ProfileControllerRemoveRequest) SetBefore(before ...jsonrpc.ClientBeforeFunc) *ProfileControllerRemoveRequest {
	remove.before = before
	return remove
}
func (remove *ProfileControllerRemoveRequest) SetAfter(after ...jsonrpc.ClientAfterFunc) *ProfileControllerRemoveRequest {
	remove.after = after
	return remove
}
func (remove *ProfileControllerRemoveRequest) WithContext(ctx context.Context) *ProfileControllerRemoveRequest {
	remove.ctx = ctx
	return remove
}
func (remove *ProfileControllerRemoveRequest) RawExecute() ([]byte, map[uint64]int, *http.Response, error) {
	return remove.c.RawExecute(remove)
}
func (remove *ProfileControllerRemoveRequest) MakeResult(data []byte) (any, error) {
	return nil, nil
}
func (remove *ProfileControllerRemoveRequest) Before() []jsonrpc.ClientBeforeFunc {
	return remove.before
}
func (remove *ProfileControllerRemoveRequest) After() []jsonrpc.ClientAfterFunc {
	return remove.after
}
func (remove *ProfileControllerRemoveRequest) Context() context.Context {
	return remove.ctx
}
func (remove *ProfileControllerRemoveRequest) Execute() (err error) {
	_, err = remove.c.Client.Execute(remove)
	if err != nil {
		return
	}
	return err
}
func (profileController *ProfileControllerClient) Create(token string, firstName string, lastName string) *ProfileControllerCreateRequest {
	r := &ProfileControllerCreateRequest{ctx: context.TODO(), c: profileController}
	r.params.token = token
	r.params.firstName = firstName
	r.params.lastName = lastName
	return r
}
func (profileController *ProfileControllerClient) Remove(id string) *ProfileControllerRemoveRequest {
	r := &ProfileControllerRemoveRequest{ctx: context.TODO(), c: profileController}
	r.params.id = id
	return r
}
func NewProfileControllerClient(target string, opts ...jsonrpc.ClientOption) *ProfileControllerClient {
	return &ProfileControllerClient{Client: jsonrpc.NewClient(target, opts...)}
}
