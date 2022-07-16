package agent

import (
	"fmt"
	"net"
	"sync"

	"github.com/mohitkumar/finch/config"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/rest"
	"github.com/mohitkumar/finch/rpc"
	"github.com/mohitkumar/finch/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Agent struct {
	Config       config.Config
	diContainer  *container.DIContiner
	httpServer   *rest.Server
	grpcServer   *grpc.Server
	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
	wg           sync.WaitGroup
}

func New(config config.Config) (*Agent, error) {
	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}
	setup := []func() error{
		a.setupDiContainer,
		a.setupHttpServer,
		a.setupGrpcServer,
	}
	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupDiContainer() error {
	a.diContainer = container.NewDiContainer()
	a.diContainer.Init(a.Config)
	return nil
}

func (a *Agent) setupHttpServer() error {
	var err error
	a.httpServer, err = rest.NewServer(a.Config.HttpPort, a.diContainer)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) setupGrpcServer() error {
	var err error
	taskService := service.NewTaskExecutionService(a.diContainer)
	conf := &rpc.GrpcConfig{
		TaskService: taskService,
	}
	a.grpcServer, err = rpc.NewGrpcServer(conf)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) Start() error {
	var err error
	a.wg.Add(2)
	go func() error {
		defer a.wg.Done()
		err = a.httpServer.Start()
		if err != nil {
			return err
		}
		return nil
	}()

	go func() error {
		defer a.wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GrpcPort))
		if err != nil {
			return err
		}
		logger.Info("startting grpc server on", zap.Int("port", a.Config.GrpcPort))

		if err := a.grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	}()
	a.wg.Wait()
	return nil
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()
	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)

	shutdown := []func() error{
		a.httpServer.Stop,
		func() error {
			a.grpcServer.GracefulStop()
			return nil
		},
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
