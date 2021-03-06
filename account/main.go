package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/amiranmanesh/go-smart-api-maker/account/proto"
	"github.com/amiranmanesh/go-smart-api-maker/account/repository"
	"github.com/amiranmanesh/go-smart-api-maker/account/server"
	"github.com/amiranmanesh/go-smart-api-maker/account/service"
	"github.com/amiranmanesh/go-smart-api-maker/utils/env"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	httpAddr = flag.String("http.addr", fmt.Sprintf(":%s", env.GetEnvItem("HTTP_PORT")), "HTTP listen address")
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "account service started")
	defer level.Info(logger).Log("msg", "account service ended")

	flag.Parse()
	ctx := context.Background()

	var srv server.IService
	{
		repository := repository.NewAccountRepository(getDataBaseModel(), logger)
		srv = service.NewService(repository, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := server.MakeEndpoint(srv)

	go func() {
		handler := server.NewHTTPServer(endpoints)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", env.GetEnvItem("GRPC_PORT")))
	if err != nil {
		panic(fmt.Sprintf("Failed to listen %v", err))
	}
	go func() {
		baseServer := grpc.NewServer()
		grpcHandler := server.NewGRPCServer(ctx, endpoints)
		proto.RegisterAccountServiceServer(baseServer, grpcHandler)
		reflection.Register(baseServer)

		errs <- baseServer.Serve(lis)
	}()

	level.Error(logger).Log("exit", <-errs)

}

func getDataBaseModel() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		env.GetEnvItem("DATABASE_USER"),
		env.GetEnvItem("DATABASE_PASSWORD"),
		env.GetEnvItem("DATABASE_HOST"),
		env.GetEnvItem("DATABASE_PORT"),
		env.GetEnvItem("DATABASE_NAME"),
	)
	connection, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		panic(errors.Trace(err))
	}
	return connection
}
