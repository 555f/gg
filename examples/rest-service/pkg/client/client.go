// Code generated by GG version dev. DO NOT EDIT.

//go:build !gg
// +build !gg

package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	dto "github.com/555f/gg/examples/rest-service/pkg/dto"
	errors "github.com/555f/gg/examples/rest-service/pkg/errors"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net"
	"net/http"
)

func instrumentRoundTripperErrCounter(counter *prometheus.CounterVec, next http.RoundTripper) promhttp.RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)
		if err != nil {
			labels := prometheus.Labels{"method": r.Method}
			errType := ""
			switch e := err.(type) {
			default:
				errType = err.Error()
			case net.Error:
				errType += "net."
				if e.Timeout() {
					errType += "timeout."
				}
				switch ee := e.(type) {
				case *net.ParseError:
					errType += "parse"
				case *net.InvalidAddrError:
					errType += "invalidAddr"
				case *net.UnknownNetworkError:
					errType += "unknownNetwork"
				case *net.DNSError:
					errType += "dns"
				case *net.OpError:
					errType += ee.Net + "." + ee.Op
				}
			}
			labels["err"] = errType
			counter.With(labels).Add(1)
		}
		return resp, err
	}
}

type outgoingInstrumentation struct {
	inflight    prometheus.Gauge
	errRequests *prometheus.CounterVec
	requests    *prometheus.CounterVec
	duration    *prometheus.HistogramVec
	dnsDuration *prometheus.HistogramVec
	tlsDuration *prometheus.HistogramVec
}

func (i *outgoingInstrumentation) Describe(in chan<- *prometheus.Desc) {
	i.inflight.Describe(in)
	i.requests.Describe(in)
	i.errRequests.Describe(in)
	i.duration.Describe(in)
	i.dnsDuration.Describe(in)
	i.tlsDuration.Describe(in)
}
func (i *outgoingInstrumentation) Collect(in chan<- prometheus.Metric) {
	i.inflight.Collect(in)
	i.requests.Collect(in)
	i.errRequests.Collect(in)
	i.duration.Collect(in)
	i.dnsDuration.Collect(in)
	i.tlsDuration.Collect(in)
}

type ClientBeforeFunc func(context.Context, *http.Request) (context.Context, error)
type ClientAfterFunc func(context.Context, *http.Response) context.Context
type clientOptions struct {
	ctx    context.Context
	before []ClientBeforeFunc
	after  []ClientAfterFunc
	client *http.Client
}
type ClientOption func(*clientOptions)

func WithContext(ctx context.Context) ClientOption {
	return func(o *clientOptions) {
		o.ctx = ctx
	}
}
func WithClient(client *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.client = client
	}
}
func WithProm(namespace string, subsystem string, reg prometheus.Registerer, constLabels map[string]string) ClientOption {
	return func(o *clientOptions) {
		i := &outgoingInstrumentation{inflight: prometheus.NewGauge(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "in_flight_requests", Help: "A gauge of in-flight outgoing requests for the client.", ConstLabels: constLabels}), requests: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "requests_total", Help: "A counter for outgoing requests from the client.", ConstLabels: constLabels}, []string{"method", "code"}), errRequests: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "err_requests_total", Help: "A counter for outgoing error requests from the client."}, []string{"method", "err"}), duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "request_duration_histogram_seconds", Help: "A histogram of outgoing request latencies.", Buckets: prometheus.DefBuckets, ConstLabels: constLabels}, []string{"method", "code"}), dnsDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "dns_duration_histogram_seconds", Help: "Trace dns latency histogram.", Buckets: prometheus.DefBuckets, ConstLabels: constLabels}, []string{"method", "code"}), tlsDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "tls_duration_histogram_seconds", Help: "Trace tls latency histogram.", Buckets: prometheus.DefBuckets, ConstLabels: constLabels}, []string{"method", "code"})}
		trace := &promhttp.InstrumentTrace{}
		o.client.Transport = instrumentRoundTripperErrCounter(i.errRequests, promhttp.InstrumentRoundTripperInFlight(i.inflight, promhttp.InstrumentRoundTripperCounter(i.requests, promhttp.InstrumentRoundTripperTrace(trace, promhttp.InstrumentRoundTripperDuration(i.duration, o.client.Transport)))))
		err := reg.Register(i)
		if err != nil {
			panic(err)
		}
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

func (r *ProfileControllerCreateRequest) SetAddress(address string) *ProfileControllerCreateRequest {
	r.params.address = &address
	return r
}
func (r *ProfileControllerClient) Create(firstName string, lastName string, address string, zip int) (profile *dto.Profile, err error) {
	profile, err = r.CreateRequest(firstName, lastName, zip).SetAddress(address).Execute()
	return
}
func (r *ProfileControllerClient) CreateRequest(firstName string, lastName string, zip int) *ProfileControllerCreateRequest {
	m := &ProfileControllerCreateRequest{client: r.opts.client, opts: &clientOptions{ctx: context.TODO()}, c: r}
	m.params.firstName = firstName
	m.params.lastName = lastName
	m.params.zip = zip
	return m
}
func (r *ProfileControllerCreateRequest) Execute(opts ...ClientOption) (profile *dto.Profile, err error) {
	for _, o := range opts {
		o(r.opts)
	}
	var body struct {
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Address   *string `json:"address,omitempty"`
		Zip       int     `json:"zip"`
	}
	ctx, cancel := context.WithCancel(r.opts.ctx)
	path := "/profiles"
	req, err := http.NewRequest("POST", r.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	body.FirstName = r.params.firstName
	body.LastName = r.params.lastName
	body.Address = r.params.address
	body.Zip = r.params.zip
	var reqData bytes.Buffer
	err = json.NewEncoder(&reqData).Encode(body)
	if err != nil {
		cancel()
		return
	}
	req.Body = io.NopCloser(&reqData)

	req.Header.Add("Content-Type", "application/json")
	before := append(r.c.opts.before, r.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := r.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(r.c.opts.after, r.opts.after...)
	for _, after := range after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		var bytes []byte
		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(bytes, &errorWrapper)
		if err != nil {
			err = fmt.Errorf("unmarshal error (%s): %w", bytes, err)
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	var respBody struct {
		Profile *dto.Profile `json:"profile"`
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	default:
		reader = resp.Body
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	}
	err = json.NewDecoder(reader).Decode(&respBody)
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

func (r *ProfileControllerClient) DownloadFile(id string) (data string, err error) {
	data, err = r.DownloadFileRequest(id).Execute()
	return
}
func (r *ProfileControllerClient) DownloadFileRequest(id string) *ProfileControllerDownloadFileRequest {
	m := &ProfileControllerDownloadFileRequest{client: r.opts.client, opts: &clientOptions{ctx: context.TODO()}, c: r}
	m.params.id = id
	return m
}
func (r *ProfileControllerDownloadFileRequest) Execute(opts ...ClientOption) (data string, err error) {
	for _, o := range opts {
		o(r.opts)
	}
	ctx, cancel := context.WithCancel(r.opts.ctx)
	path := fmt.Sprintf("/profiles/%s/file", r.params.id)
	req, err := http.NewRequest("GET", r.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	req.Header.Add("Content-Type", "application/json")
	before := append(r.c.opts.before, r.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := r.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(r.c.opts.after, r.opts.after...)
	for _, after := range after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		var bytes []byte
		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(bytes, &errorWrapper)
		if err != nil {
			err = fmt.Errorf("unmarshal error (%s): %w", bytes, err)
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	var respBody struct {
		Data string `json:"data"`
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	default:
		reader = resp.Body
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	}
	err = json.NewDecoder(reader).Decode(&respBody)
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

func (r *ProfileControllerClient) Remove(id string) (err error) {
	err = r.RemoveRequest(id).Execute()
	return
}
func (r *ProfileControllerClient) RemoveRequest(id string) *ProfileControllerRemoveRequest {
	m := &ProfileControllerRemoveRequest{client: r.opts.client, opts: &clientOptions{ctx: context.TODO()}, c: r}
	m.params.id = id
	return m
}
func (r *ProfileControllerRemoveRequest) Execute(opts ...ClientOption) (err error) {
	for _, o := range opts {
		o(r.opts)
	}
	var body struct {
		Id string `json:"id"`
	}
	ctx, cancel := context.WithCancel(r.opts.ctx)
	path := "/profiles/{id}"
	req, err := http.NewRequest("DELETE", r.c.target+path, nil)
	if err != nil {
		cancel()
		return
	}

	body.Id = r.params.id
	var reqData bytes.Buffer
	err = json.NewEncoder(&reqData).Encode(body)
	if err != nil {
		cancel()
		return
	}
	req.Body = io.NopCloser(&reqData)

	req.Header.Add("Content-Type", "application/json")
	before := append(r.c.opts.before, r.opts.before...)
	for _, before := range before {
		ctx, err = before(ctx, req)
		if err != nil {
			cancel()
			return
		}
	}
	resp, err := r.client.Do(req)
	if err != nil {
		cancel()
		return
	}
	after := append(r.c.opts.after, r.opts.after...)
	for _, after := range after {
		ctx = after(ctx, resp)
	}
	defer resp.Body.Close()
	defer cancel()
	if resp.StatusCode > 399 {
		var errorWrapper errors.ErrorWrapper
		var bytes []byte
		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(bytes, &errorWrapper)
		if err != nil {
			err = fmt.Errorf("unmarshal error (%s): %w", bytes, err)
			return
		}
		err = &errors.DefaultError{Data: errorWrapper.Data, ErrorText: errorWrapper.ErrorText, Code: errorWrapper.Code}
		return
	}
	return nil
}
func NewProfileControllerClient(target string, opts ...ClientOption) *ProfileControllerClient {
	c := &ProfileControllerClient{target: target, opts: &clientOptions{client: http.DefaultClient}}
	for _, o := range opts {
		o(c.opts)
	}
	return c
}
