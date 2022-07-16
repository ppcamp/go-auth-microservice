package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/ppcamp/go-auth/src/configs"
	handlers "github.com/ppcamp/go-auth/src/http"
	"github.com/ppcamp/go-auth/src/http/gRPC/auth"
	"github.com/ppcamp/go-auth/src/http/gRPC/user_password"
	"github.com/ppcamp/go-auth/src/repositories/cache"
	"github.com/ppcamp/go-auth/src/repositories/database"
	grpcutils "github.com/ppcamp/go-auth/src/utils/grpc"
	"github.com/ppcamp/go-auth/src/utils/jwt"

	"github.com/ppcamp/go-cli/env"
	"github.com/ppcamp/go-cli/shutdown"
	xtenderrors "github.com/ppcamp/go-xtendlib/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcServer *grpc.Server
	server     *http.Server
	handler    *handlers.Handler
)

// load flags from the environment and assign the values to each variable in configs pkg
func init() {
	flags := configs.Flags()
	err := env.Parse(flags)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	ctx := context.Background()

	before()

	err := shutdown.Graceful(ctx, run)
	if err != nil {
		log.Panic(err)
	}

	after()
}

// before is usually used to initialize the services
func before() {
	// define cache config
	cacheAddr := fmt.Sprintf(cache.CONNECTION_STRING, configs.CacheHost, configs.CachePort)
	cacheId := fmt.Sprintf(cache.ID_STRING_FORMAT, configs.APP_NAME, configs.AppId)
	cacheConfig := cache.CacheConfig{
		Addr:        cacheAddr,
		Password:    configs.CachePassword, // no password set
		DB:          configs.CacheDb,       // use default DB
		DialTimeout: configs.CACHE_CONNECTION_TIMEOUT,
	}

	// initialize cacheRepository
	log.WithFields(
		log.Fields{"Id": cacheId, "Config": cacheConfig}).Info("Starting connection with cache")
	cacheRepository := xtenderrors.Must(cache.NewCacheRepository(cacheConfig, cacheId))

	// initialize jwt vault manager
	log.Info("Creating vault manager/signer")
	privateKey := xtenderrors.Must(jwt.ParseSSHPrivateKey(configs.JwtPrivate))
	signer := jwt.DefaultSigner(privateKey)

	// define database config
	connQuery := fmt.Sprintf(
		database.CONNECTION_QUERY,
		configs.DatabaseHost,
		configs.DatabasePort,
		configs.DatabaseUser,
		configs.DatabasePassword,
		configs.DatabaseName,
	)

	// initialize database
	log.WithField("ConnectionQuery", connQuery).Info("Starting a new store")
	db := xtenderrors.Must(database.NewStore(connQuery))

	log.Info("Initializing handlers")
	handler = &handlers.Handler{
		Cache:    cacheRepository,
		Database: db,
		Signer:   signer,
	}
}

// after is called after the server closed, usually is a clean up function
func after() {
	log.Info("Closing server...")
	defer log.Info("Server closed!")

	log.Info("\t - Closing gRPC connections")
	grpcServer.GracefulStop()

	log.Info("\t - Closing http server")
	if err := server.Close(); err != nil {
		log.Error(err)
	}

	log.Info("\t - Closing Database connections")
	if err := handler.Database.Close(); err != nil {
		log.Error("fail when closing database")
	}

	log.Info("\t - Closing cache connections")
	if err := handler.Cache.Close(); err != nil {
		log.Error(err)
	}
}

func run(ctx context.Context) error {
	grpcServer = grpc.NewServer()

	// initialize gRPC services
	authServer := auth.NewAuthService(handler)
	userServer := user_password.NewUserPasswordService(handler)
	health := &grpcutils.GrpcHealthService{}

	// register services in gRPC
	grpc_health_v1.RegisterHealthServer(grpcServer, health)
	auth.RegisterAuthServiceServer(grpcServer, authServer)
	user_password.RegisterUserPasswordServiceServer(grpcServer, userServer)

	// initialize tcp listener
	listener, err := net.Listen("tcp", configs.AppPort)
	if err != nil {
		return err
	}

	// make gRPC http server
	server = grpcutils.NewMuxServer(http.NewServeMux(), grpcServer)

	//start server
	log.Infof("Server listening at http://localhost%s", configs.AppPort)

	err = server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return xtenderrors.Wraps("fail to serve gRPC", err)
	}

	return nil
}
