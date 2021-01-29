package fileservice

import (
	"compress/flate"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	envconfig "github.com/unistack-org/micro-config-env/v3"
	fileconfig "github.com/unistack-org/micro-config-file/v3"
	promwrapper "github.com/unistack-org/micro-metrics-prometheus/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	s3store "github.com/unistack-org/micro-store-s3/v3"
	idwrapper "github.com/unistack-org/micro-wrapper-requestid/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/store"
	"github.com/vielendanke/file-service/configs"
	"github.com/vielendanke/file-service/internal/app/fileservice/commons/http/middleware"
	"github.com/vielendanke/file-service/internal/app/fileservice/commons/stats"
	"github.com/vielendanke/file-service/internal/app/fileservice/handlers"
	"github.com/vielendanke/file-service/internal/app/fileservice/middlewares"
	"github.com/vielendanke/file-service/internal/app/fileservice/repository"
	"github.com/vielendanke/file-service/internal/app/fileservice/service"
	pb "github.com/vielendanke/file-service/proto"
)

func initDB(name, url string, errCh chan<- error) <-chan *sqlx.DB {
	dbCh := make(chan *sqlx.DB, 1)
	go func() {
		defer close(dbCh)
		db, err := sqlx.Connect(name, url)
		if err != nil {
			errCh <- err
		}
		dbCh <- db
	}()
	return dbCh
}

// StartFileService ...
func StartFileService(ctx context.Context, errCh chan<- error) {
	cfg := configs.NewConfig("file-service", "1.0")

	if err := config.Load(ctx,
		config.NewConfig( // load from defaults
			config.Struct(cfg), // pass config struct
		),
		fileconfig.NewConfig( // load from file
			config.AllowFail(true),             // that may be not exists
			config.Struct(cfg),                 // pass config struct
			config.Codec(jsoncodec.NewCodec()), // file config in json
			fileconfig.Path("./local.json"),    // nearby file
		),
		envconfig.NewConfig( // load from environment
			config.Struct(cfg), // pass config struct
		),
	); err != nil {
		errCh <- err
	}

	s3Dirty := s3store.NewStore(
		store.Name("dirty_region"),
		s3store.AccessKey(cfg.Amazon.DirtyRegion.AccessKey),
		s3store.SecretKey(cfg.Amazon.DirtyRegion.SecretKey),
		s3store.Endpoint(cfg.Amazon.DirtyRegion.Endpoint),
	)
	s3Clean := s3store.NewStore(
		store.Name("clean_region"),
		s3store.AccessKey(cfg.Amazon.CleanRegion.AccessKey),
		s3store.SecretKey(cfg.Amazon.CleanRegion.SecretKey),
		s3store.Endpoint(cfg.Amazon.CleanRegion.Endpoint),
	)
	if err := s3Dirty.Init(); err != nil {
		errCh <- err
	}
	if err := s3Dirty.Connect(ctx); err != nil {
		errCh <- err
	}
	defer func() {
		err := s3Dirty.Disconnect(ctx)
		if err != nil {
			logger.Errorf(ctx, "Error during disconnect from s3, %v", err)
		}
	}()
	if err := s3Clean.Init(); err != nil {
		errCh <- err
	}
	if err := s3Clean.Connect(ctx); err != nil {
		errCh <- err
	}
	defer func() {
		err := s3Clean.Disconnect(ctx)
		if err != nil {
			logger.Errorf(ctx, "Error during disconnect from s3, %v", err)
		}
	}()

	options := append([]micro.Option{},
		micro.Servers(httpsrv.NewServer()),
		micro.Context(ctx),
		micro.Name(cfg.Server.Name),
		micro.Version(cfg.Server.Version),
		micro.Stores(s3Clean, s3Dirty),
	)
	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		errCh <- err
	}

	if err := svc.Init(
		micro.Servers(httpsrv.NewServer(
			server.Name(cfg.Server.Name),
			server.Version(cfg.Server.Version),
			server.Address(cfg.Server.Addr),
			server.Context(ctx),
			server.Codec("application/json", jsoncodec.NewCodec()),
			server.WrapHandler(promwrapper.NewHandlerWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
			server.WrapHandler(idwrapper.NewServerHandlerWrapper()),
		)),
		micro.Clients(httpcli.NewClient(
			client.ContentType("application/json"),
			client.Codec("application/json", jsoncodec.NewCodec()),
			client.Wrap(promwrapper.NewClientWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
			client.Wrap(idwrapper.NewClientWrapper()),
		)),
		micro.Stores(s3Clean, s3Dirty),
	); err != nil {
		errCh <- err
	}
	router := mux.NewRouter()

	ctm := middlewares.NewContentTypeMiddleware("application/json")

	router.Use(ctm.ContentTypeMiddleware)
	router.Use(middleware.HttpMetricsWrapper)
	router.Use(middleware.NewRequestIDMiddleware().Wrapper)
	router.Use(middleware.NewLoggerMiddleware().Wrapper)
	router.Use(middleware.NewNocacheMiddleware().Wrapper)
	router.Use(middleware.NewCompressMiddleware(flate.BestSpeed).Wrapper)

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(ctx, "Not found, %v/n", r.URL)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("Not found. Path: %s, Method: %s", r.RequestURI, r.Method)))
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(ctx, "Method not allowed, %v/n", r.URL)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(fmt.Sprintf("Method not allowed, %s", r.Method)))
	})
	endpoints := pb.NewFileProcessingServiceEndpoints()

	db := <-initDB("postgres", cfg.Database.URL, errCh)

	fr := repository.NewAWSFileRepository(db)

	srv := service.NewAWSProcessingService(jsoncodec.NewCodec(), fr, svc.Store("clean_region"), svc.Store("dirty_region"))

	handler := handlers.NewFileServiceHandler(srv, jsoncodec.NewCodec())

	if err := configs.ConfigureHandlerToEndpoints(router, handler, endpoints); err != nil {
		errCh <- err
	}
	if err := svc.Server().Handle(svc.Server().NewHandler(router)); err != nil {
		errCh <- err
	}

	statsOpts := append([]stats.Option{},
		stats.WithDefaultHealth(),
		stats.WithMetrics(),
		stats.WithVersionDate(cfg.Server.Name, time.Now().String()),
	)

	healthServer := stats.NewServer(statsOpts...)
	go func() {
		logger.Fatal(ctx, healthServer.Serve(cfg.Metric.Addr))
	}()

	if err := svc.Run(); err != nil {
		errCh <- err
	}
}
