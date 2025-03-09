// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: daemon/daemon.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DaemonService_Login_FullMethodName                    = "/daemon.DaemonService/Login"
	DaemonService_WaitSSOLogin_FullMethodName             = "/daemon.DaemonService/WaitSSOLogin"
	DaemonService_Up_FullMethodName                       = "/daemon.DaemonService/Up"
	DaemonService_Status_FullMethodName                   = "/daemon.DaemonService/Status"
	DaemonService_Down_FullMethodName                     = "/daemon.DaemonService/Down"
	DaemonService_GetConfig_FullMethodName                = "/daemon.DaemonService/GetConfig"
	DaemonService_ListNetworks_FullMethodName             = "/daemon.DaemonService/ListNetworks"
	DaemonService_SelectNetworks_FullMethodName           = "/daemon.DaemonService/SelectNetworks"
	DaemonService_DeselectNetworks_FullMethodName         = "/daemon.DaemonService/DeselectNetworks"
	DaemonService_ForwardingRules_FullMethodName          = "/daemon.DaemonService/ForwardingRules"
	DaemonService_DebugBundle_FullMethodName              = "/daemon.DaemonService/DebugBundle"
	DaemonService_GetLogLevel_FullMethodName              = "/daemon.DaemonService/GetLogLevel"
	DaemonService_SetLogLevel_FullMethodName              = "/daemon.DaemonService/SetLogLevel"
	DaemonService_ListStates_FullMethodName               = "/daemon.DaemonService/ListStates"
	DaemonService_CleanState_FullMethodName               = "/daemon.DaemonService/CleanState"
	DaemonService_DeleteState_FullMethodName              = "/daemon.DaemonService/DeleteState"
	DaemonService_SetNetworkMapPersistence_FullMethodName = "/daemon.DaemonService/SetNetworkMapPersistence"
	DaemonService_TracePacket_FullMethodName              = "/daemon.DaemonService/TracePacket"
	DaemonService_SubscribeEvents_FullMethodName          = "/daemon.DaemonService/SubscribeEvents"
	DaemonService_GetEvents_FullMethodName                = "/daemon.DaemonService/GetEvents"
)

// DaemonServiceClient is the client API for DaemonService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaemonServiceClient interface {
	// Login uses setup key to prepare configuration for the daemon.
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	// WaitSSOLogin uses the userCode to validate the TokenInfo and
	// waits for the user to continue with the login on a browser
	WaitSSOLogin(ctx context.Context, in *WaitSSOLoginRequest, opts ...grpc.CallOption) (*WaitSSOLoginResponse, error)
	// Up starts engine work in the daemon.
	Up(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*UpResponse, error)
	// Status of the service.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// Down engine work in the daemon.
	Down(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (*DownResponse, error)
	// GetConfig of the daemon.
	GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	// List available networks
	ListNetworks(ctx context.Context, in *ListNetworksRequest, opts ...grpc.CallOption) (*ListNetworksResponse, error)
	// Select specific routes
	SelectNetworks(ctx context.Context, in *SelectNetworksRequest, opts ...grpc.CallOption) (*SelectNetworksResponse, error)
	// Deselect specific routes
	DeselectNetworks(ctx context.Context, in *SelectNetworksRequest, opts ...grpc.CallOption) (*SelectNetworksResponse, error)
	ForwardingRules(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ForwardingRulesResponse, error)
	// DebugBundle creates a debug bundle
	DebugBundle(ctx context.Context, in *DebugBundleRequest, opts ...grpc.CallOption) (*DebugBundleResponse, error)
	// GetLogLevel gets the log level of the daemon
	GetLogLevel(ctx context.Context, in *GetLogLevelRequest, opts ...grpc.CallOption) (*GetLogLevelResponse, error)
	// SetLogLevel sets the log level of the daemon
	SetLogLevel(ctx context.Context, in *SetLogLevelRequest, opts ...grpc.CallOption) (*SetLogLevelResponse, error)
	// List all states
	ListStates(ctx context.Context, in *ListStatesRequest, opts ...grpc.CallOption) (*ListStatesResponse, error)
	// Clean specific state or all states
	CleanState(ctx context.Context, in *CleanStateRequest, opts ...grpc.CallOption) (*CleanStateResponse, error)
	// Delete specific state or all states
	DeleteState(ctx context.Context, in *DeleteStateRequest, opts ...grpc.CallOption) (*DeleteStateResponse, error)
	// SetNetworkMapPersistence enables or disables network map persistence
	SetNetworkMapPersistence(ctx context.Context, in *SetNetworkMapPersistenceRequest, opts ...grpc.CallOption) (*SetNetworkMapPersistenceResponse, error)
	TracePacket(ctx context.Context, in *TracePacketRequest, opts ...grpc.CallOption) (*TracePacketResponse, error)
	SubscribeEvents(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SystemEvent], error)
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
}

type daemonServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDaemonServiceClient(cc grpc.ClientConnInterface) DaemonServiceClient {
	return &daemonServiceClient{cc}
}

func (c *daemonServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, DaemonService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) WaitSSOLogin(ctx context.Context, in *WaitSSOLoginRequest, opts ...grpc.CallOption) (*WaitSSOLoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WaitSSOLoginResponse)
	err := c.cc.Invoke(ctx, DaemonService_WaitSSOLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Up(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*UpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpResponse)
	err := c.cc.Invoke(ctx, DaemonService_Up_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, DaemonService_Status_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Down(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (*DownResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownResponse)
	err := c.cc.Invoke(ctx, DaemonService_Down_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, DaemonService_GetConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) ListNetworks(ctx context.Context, in *ListNetworksRequest, opts ...grpc.CallOption) (*ListNetworksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListNetworksResponse)
	err := c.cc.Invoke(ctx, DaemonService_ListNetworks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) SelectNetworks(ctx context.Context, in *SelectNetworksRequest, opts ...grpc.CallOption) (*SelectNetworksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SelectNetworksResponse)
	err := c.cc.Invoke(ctx, DaemonService_SelectNetworks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) DeselectNetworks(ctx context.Context, in *SelectNetworksRequest, opts ...grpc.CallOption) (*SelectNetworksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SelectNetworksResponse)
	err := c.cc.Invoke(ctx, DaemonService_DeselectNetworks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) ForwardingRules(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ForwardingRulesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ForwardingRulesResponse)
	err := c.cc.Invoke(ctx, DaemonService_ForwardingRules_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) DebugBundle(ctx context.Context, in *DebugBundleRequest, opts ...grpc.CallOption) (*DebugBundleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DebugBundleResponse)
	err := c.cc.Invoke(ctx, DaemonService_DebugBundle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) GetLogLevel(ctx context.Context, in *GetLogLevelRequest, opts ...grpc.CallOption) (*GetLogLevelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetLogLevelResponse)
	err := c.cc.Invoke(ctx, DaemonService_GetLogLevel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) SetLogLevel(ctx context.Context, in *SetLogLevelRequest, opts ...grpc.CallOption) (*SetLogLevelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetLogLevelResponse)
	err := c.cc.Invoke(ctx, DaemonService_SetLogLevel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) ListStates(ctx context.Context, in *ListStatesRequest, opts ...grpc.CallOption) (*ListStatesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListStatesResponse)
	err := c.cc.Invoke(ctx, DaemonService_ListStates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) CleanState(ctx context.Context, in *CleanStateRequest, opts ...grpc.CallOption) (*CleanStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CleanStateResponse)
	err := c.cc.Invoke(ctx, DaemonService_CleanState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) DeleteState(ctx context.Context, in *DeleteStateRequest, opts ...grpc.CallOption) (*DeleteStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteStateResponse)
	err := c.cc.Invoke(ctx, DaemonService_DeleteState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) SetNetworkMapPersistence(ctx context.Context, in *SetNetworkMapPersistenceRequest, opts ...grpc.CallOption) (*SetNetworkMapPersistenceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetNetworkMapPersistenceResponse)
	err := c.cc.Invoke(ctx, DaemonService_SetNetworkMapPersistence_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) TracePacket(ctx context.Context, in *TracePacketRequest, opts ...grpc.CallOption) (*TracePacketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TracePacketResponse)
	err := c.cc.Invoke(ctx, DaemonService_TracePacket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) SubscribeEvents(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SystemEvent], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DaemonService_ServiceDesc.Streams[0], DaemonService_SubscribeEvents_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeRequest, SystemEvent]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DaemonService_SubscribeEventsClient = grpc.ServerStreamingClient[SystemEvent]

func (c *daemonServiceClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEventsResponse)
	err := c.cc.Invoke(ctx, DaemonService_GetEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServiceServer is the server API for DaemonService service.
// All implementations must embed UnimplementedDaemonServiceServer
// for forward compatibility.
type DaemonServiceServer interface {
	// Login uses setup key to prepare configuration for the daemon.
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// WaitSSOLogin uses the userCode to validate the TokenInfo and
	// waits for the user to continue with the login on a browser
	WaitSSOLogin(context.Context, *WaitSSOLoginRequest) (*WaitSSOLoginResponse, error)
	// Up starts engine work in the daemon.
	Up(context.Context, *UpRequest) (*UpResponse, error)
	// Status of the service.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	// Down engine work in the daemon.
	Down(context.Context, *DownRequest) (*DownResponse, error)
	// GetConfig of the daemon.
	GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	// List available networks
	ListNetworks(context.Context, *ListNetworksRequest) (*ListNetworksResponse, error)
	// Select specific routes
	SelectNetworks(context.Context, *SelectNetworksRequest) (*SelectNetworksResponse, error)
	// Deselect specific routes
	DeselectNetworks(context.Context, *SelectNetworksRequest) (*SelectNetworksResponse, error)
	ForwardingRules(context.Context, *EmptyRequest) (*ForwardingRulesResponse, error)
	// DebugBundle creates a debug bundle
	DebugBundle(context.Context, *DebugBundleRequest) (*DebugBundleResponse, error)
	// GetLogLevel gets the log level of the daemon
	GetLogLevel(context.Context, *GetLogLevelRequest) (*GetLogLevelResponse, error)
	// SetLogLevel sets the log level of the daemon
	SetLogLevel(context.Context, *SetLogLevelRequest) (*SetLogLevelResponse, error)
	// List all states
	ListStates(context.Context, *ListStatesRequest) (*ListStatesResponse, error)
	// Clean specific state or all states
	CleanState(context.Context, *CleanStateRequest) (*CleanStateResponse, error)
	// Delete specific state or all states
	DeleteState(context.Context, *DeleteStateRequest) (*DeleteStateResponse, error)
	// SetNetworkMapPersistence enables or disables network map persistence
	SetNetworkMapPersistence(context.Context, *SetNetworkMapPersistenceRequest) (*SetNetworkMapPersistenceResponse, error)
	TracePacket(context.Context, *TracePacketRequest) (*TracePacketResponse, error)
	SubscribeEvents(*SubscribeRequest, grpc.ServerStreamingServer[SystemEvent]) error
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
	mustEmbedUnimplementedDaemonServiceServer()
}

// UnimplementedDaemonServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDaemonServiceServer struct{}

func (UnimplementedDaemonServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedDaemonServiceServer) WaitSSOLogin(context.Context, *WaitSSOLoginRequest) (*WaitSSOLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitSSOLogin not implemented")
}
func (UnimplementedDaemonServiceServer) Up(context.Context, *UpRequest) (*UpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Up not implemented")
}
func (UnimplementedDaemonServiceServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedDaemonServiceServer) Down(context.Context, *DownRequest) (*DownResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Down not implemented")
}
func (UnimplementedDaemonServiceServer) GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedDaemonServiceServer) ListNetworks(context.Context, *ListNetworksRequest) (*ListNetworksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNetworks not implemented")
}
func (UnimplementedDaemonServiceServer) SelectNetworks(context.Context, *SelectNetworksRequest) (*SelectNetworksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectNetworks not implemented")
}
func (UnimplementedDaemonServiceServer) DeselectNetworks(context.Context, *SelectNetworksRequest) (*SelectNetworksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeselectNetworks not implemented")
}
func (UnimplementedDaemonServiceServer) ForwardingRules(context.Context, *EmptyRequest) (*ForwardingRulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardingRules not implemented")
}
func (UnimplementedDaemonServiceServer) DebugBundle(context.Context, *DebugBundleRequest) (*DebugBundleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DebugBundle not implemented")
}
func (UnimplementedDaemonServiceServer) GetLogLevel(context.Context, *GetLogLevelRequest) (*GetLogLevelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogLevel not implemented")
}
func (UnimplementedDaemonServiceServer) SetLogLevel(context.Context, *SetLogLevelRequest) (*SetLogLevelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLogLevel not implemented")
}
func (UnimplementedDaemonServiceServer) ListStates(context.Context, *ListStatesRequest) (*ListStatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStates not implemented")
}
func (UnimplementedDaemonServiceServer) CleanState(context.Context, *CleanStateRequest) (*CleanStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanState not implemented")
}
func (UnimplementedDaemonServiceServer) DeleteState(context.Context, *DeleteStateRequest) (*DeleteStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteState not implemented")
}
func (UnimplementedDaemonServiceServer) SetNetworkMapPersistence(context.Context, *SetNetworkMapPersistenceRequest) (*SetNetworkMapPersistenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNetworkMapPersistence not implemented")
}
func (UnimplementedDaemonServiceServer) TracePacket(context.Context, *TracePacketRequest) (*TracePacketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TracePacket not implemented")
}
func (UnimplementedDaemonServiceServer) SubscribeEvents(*SubscribeRequest, grpc.ServerStreamingServer[SystemEvent]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeEvents not implemented")
}
func (UnimplementedDaemonServiceServer) GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (UnimplementedDaemonServiceServer) mustEmbedUnimplementedDaemonServiceServer() {}
func (UnimplementedDaemonServiceServer) testEmbeddedByValue()                       {}

// UnsafeDaemonServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaemonServiceServer will
// result in compilation errors.
type UnsafeDaemonServiceServer interface {
	mustEmbedUnimplementedDaemonServiceServer()
}

func RegisterDaemonServiceServer(s grpc.ServiceRegistrar, srv DaemonServiceServer) {
	// If the following call pancis, it indicates UnimplementedDaemonServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DaemonService_ServiceDesc, srv)
}

func _DaemonService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_WaitSSOLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WaitSSOLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).WaitSSOLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_WaitSSOLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).WaitSSOLogin(ctx, req.(*WaitSSOLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Up_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Up(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_Up_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Up(ctx, req.(*UpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Down_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Down(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_Down_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Down(ctx, req.(*DownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_GetConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).GetConfig(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_ListNetworks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNetworksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).ListNetworks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_ListNetworks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).ListNetworks(ctx, req.(*ListNetworksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_SelectNetworks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectNetworksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).SelectNetworks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_SelectNetworks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).SelectNetworks(ctx, req.(*SelectNetworksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_DeselectNetworks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectNetworksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).DeselectNetworks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_DeselectNetworks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).DeselectNetworks(ctx, req.(*SelectNetworksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_ForwardingRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).ForwardingRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_ForwardingRules_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).ForwardingRules(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_DebugBundle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DebugBundleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).DebugBundle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_DebugBundle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).DebugBundle(ctx, req.(*DebugBundleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_GetLogLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogLevelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).GetLogLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_GetLogLevel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).GetLogLevel(ctx, req.(*GetLogLevelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_SetLogLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetLogLevelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).SetLogLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_SetLogLevel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).SetLogLevel(ctx, req.(*SetLogLevelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_ListStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).ListStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_ListStates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).ListStates(ctx, req.(*ListStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_CleanState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CleanStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).CleanState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_CleanState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).CleanState(ctx, req.(*CleanStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_DeleteState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).DeleteState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_DeleteState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).DeleteState(ctx, req.(*DeleteStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_SetNetworkMapPersistence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetNetworkMapPersistenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).SetNetworkMapPersistence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_SetNetworkMapPersistence_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).SetNetworkMapPersistence(ctx, req.(*SetNetworkMapPersistenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_TracePacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TracePacketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).TracePacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_TracePacket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).TracePacket(ctx, req.(*TracePacketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_SubscribeEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DaemonServiceServer).SubscribeEvents(m, &grpc.GenericServerStream[SubscribeRequest, SystemEvent]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DaemonService_SubscribeEventsServer = grpc.ServerStreamingServer[SystemEvent]

func _DaemonService_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaemonService_GetEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DaemonService_ServiceDesc is the grpc.ServiceDesc for DaemonService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DaemonService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "daemon.DaemonService",
	HandlerType: (*DaemonServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _DaemonService_Login_Handler,
		},
		{
			MethodName: "WaitSSOLogin",
			Handler:    _DaemonService_WaitSSOLogin_Handler,
		},
		{
			MethodName: "Up",
			Handler:    _DaemonService_Up_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _DaemonService_Status_Handler,
		},
		{
			MethodName: "Down",
			Handler:    _DaemonService_Down_Handler,
		},
		{
			MethodName: "GetConfig",
			Handler:    _DaemonService_GetConfig_Handler,
		},
		{
			MethodName: "ListNetworks",
			Handler:    _DaemonService_ListNetworks_Handler,
		},
		{
			MethodName: "SelectNetworks",
			Handler:    _DaemonService_SelectNetworks_Handler,
		},
		{
			MethodName: "DeselectNetworks",
			Handler:    _DaemonService_DeselectNetworks_Handler,
		},
		{
			MethodName: "ForwardingRules",
			Handler:    _DaemonService_ForwardingRules_Handler,
		},
		{
			MethodName: "DebugBundle",
			Handler:    _DaemonService_DebugBundle_Handler,
		},
		{
			MethodName: "GetLogLevel",
			Handler:    _DaemonService_GetLogLevel_Handler,
		},
		{
			MethodName: "SetLogLevel",
			Handler:    _DaemonService_SetLogLevel_Handler,
		},
		{
			MethodName: "ListStates",
			Handler:    _DaemonService_ListStates_Handler,
		},
		{
			MethodName: "CleanState",
			Handler:    _DaemonService_CleanState_Handler,
		},
		{
			MethodName: "DeleteState",
			Handler:    _DaemonService_DeleteState_Handler,
		},
		{
			MethodName: "SetNetworkMapPersistence",
			Handler:    _DaemonService_SetNetworkMapPersistence_Handler,
		},
		{
			MethodName: "TracePacket",
			Handler:    _DaemonService_TracePacket_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _DaemonService_GetEvents_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeEvents",
			Handler:       _DaemonService_SubscribeEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "daemon/daemon.proto",
}
