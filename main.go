package main

import (
	"log"
	"net"
	"os"
	"path"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	socket := path.Join(pluginapi.DevicePluginPath, "virtual_device.sock")
	os.Remove(socket)
	go func() {
		// ping server
		conn, err := dial(socket, 5*time.Second)
		if err != nil {
			panic(err)
		}
		conn.Close()
		log.Println("ping server success")
		// register to kubelet after server started
		conn, err = dial(pluginapi.KubeletSocket, 5*time.Second)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		client := pluginapi.NewRegistrationClient(conn)
		reqt := &pluginapi.RegisterRequest{
			Version:      pluginapi.Version,
			Endpoint:     path.Base(socket),
			ResourceName: "device.plugin/device",
			Options:      &pluginapi.DevicePluginOptions{},
		}
		_, err = client.Register(context.Background(), reqt)
		if err != nil {
			panic(err)
		}
		log.Println("register server success")
	}()
	// startup server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	sock, err := net.Listen("unix", socket)
	if err != nil {
		panic(err)
	}
	pluginapi.RegisterDevicePluginServer(server, &Service{})
	log.Println("running server")
	panic(server.Serve(sock))
}

type Service struct{}

// GetDevicePluginOptions returns the values of the optional settings for this plugin
func (m *Service) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	log.Println("GetDevicePluginOptions")
	return &pluginapi.DevicePluginOptions{}, nil
}

// ListAndWatch lists devices and update that list according to the health status
func (m *Service) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{
		Devices: []*pluginapi.Device{
			{
				ID:     "0",
				Health: pluginapi.Healthy,
			}, {
				ID:     "1",
				Health: pluginapi.Healthy,
			},
		},
	})
	log.Println("send devices")
	select {}
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (m *Service) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	panic("unimplemented")
}

// Allocate which return list of devices.
func (m *Service) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	res := &pluginapi.AllocateResponse{}
	for _, r := range reqs.GetContainerRequests() {
		devices := &pluginapi.ContainerAllocateResponse{}
		devices.Envs = make(map[string]string)
		for _, id := range r.GetDevicesIDs() {
			devices.Envs["DEVICE_ID"+id] = id
		}
		res.ContainerResponses = append(res.ContainerResponses, devices)
	}
	log.Println("allocate devices")
	return res, nil
}

// PreStartContainer is unimplemented for this plugin
func (m *Service) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	log.Println("PreStartContainer")
	return &pluginapi.PreStartContainerResponse{}, nil
}

// dial establishes the gRPC communication with the registered device plugin.
func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
