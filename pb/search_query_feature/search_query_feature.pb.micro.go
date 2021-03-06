// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: search_query_feature.proto

package search_query_feature

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Feature service

type FeatureService interface {
	Get(ctx context.Context, in *OnlineRequest, opts ...client.CallOption) (*OnlineResponse, error)
	Set(ctx context.Context, in *OfflineRequest, opts ...client.CallOption) (*OfflineResponse, error)
	Chose(ctx context.Context, in *ChoseRequest, opts ...client.CallOption) (*ChoseResponse, error)
}

type featureService struct {
	c    client.Client
	name string
}

func NewFeatureService(name string, c client.Client) FeatureService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "search_query_feature"
	}
	return &featureService{
		c:    c,
		name: name,
	}
}

func (c *featureService) Get(ctx context.Context, in *OnlineRequest, opts ...client.CallOption) (*OnlineResponse, error) {
	req := c.c.NewRequest(c.name, "Feature.Get", in)
	out := new(OnlineResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *featureService) Set(ctx context.Context, in *OfflineRequest, opts ...client.CallOption) (*OfflineResponse, error) {
	req := c.c.NewRequest(c.name, "Feature.Set", in)
	out := new(OfflineResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *featureService) Chose(ctx context.Context, in *ChoseRequest, opts ...client.CallOption) (*ChoseResponse, error) {
	req := c.c.NewRequest(c.name, "Feature.Chose", in)
	out := new(ChoseResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Feature service

type FeatureHandler interface {
	Get(context.Context, *OnlineRequest, *OnlineResponse) error
	Set(context.Context, *OfflineRequest, *OfflineResponse) error
	Chose(context.Context, *ChoseRequest, *ChoseResponse) error
}

func RegisterFeatureHandler(s server.Server, hdlr FeatureHandler, opts ...server.HandlerOption) error {
	type feature interface {
		Get(ctx context.Context, in *OnlineRequest, out *OnlineResponse) error
		Set(ctx context.Context, in *OfflineRequest, out *OfflineResponse) error
		Chose(ctx context.Context, in *ChoseRequest, out *ChoseResponse) error
	}
	type Feature struct {
		feature
	}
	h := &featureHandler{hdlr}
	return s.Handle(s.NewHandler(&Feature{h}, opts...))
}

type featureHandler struct {
	FeatureHandler
}

func (h *featureHandler) Get(ctx context.Context, in *OnlineRequest, out *OnlineResponse) error {
	return h.FeatureHandler.Get(ctx, in, out)
}

func (h *featureHandler) Set(ctx context.Context, in *OfflineRequest, out *OfflineResponse) error {
	return h.FeatureHandler.Set(ctx, in, out)
}

func (h *featureHandler) Chose(ctx context.Context, in *ChoseRequest, out *ChoseResponse) error {
	return h.FeatureHandler.Chose(ctx, in, out)
}
