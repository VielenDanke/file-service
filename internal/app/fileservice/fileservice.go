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
	promwrapper "github.com/unistack-org/micro-metrics-prometheus/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	s3store "github.com/unistack-org/micro-store-s3/v3"
	idwrapper "github.com/unistack-org/micro-wrapper-requestid/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	"github.com/vielendanke/commons/http/middleware"
	"github.com/vielendanke/commons/stats"
	"github.com/vielendanke/file-service/configs"
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
	s3 := s3store.NewStore(
		s3store.AccessKey("Q3AM3UQ867SPQQA43P2F"),
		s3store.SecretKey("zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"),
		s3store.Endpoint("https://play.minio.io"),
	)
	if err := s3.Init(); err != nil {
		errCh <- err
	}
	if err := s3.Connect(ctx); err != nil {
		errCh <- err
	}
	defer func() {
		err := s3.Disconnect(ctx)
		if err != nil {
			logger.Errorf(ctx, "Error during disconnect from s3, %v", err)
		}
	}()

	if err := config.Load(ctx); err != nil {
		errCh <- err
	}

	options := append([]micro.Option{},
		micro.Server(httpsrv.NewServer()),
		micro.Client(httpcli.NewClient()),
		micro.Context(ctx),
		micro.Name("file-service"),
		micro.Version("1.0"),
		micro.Store(s3),
	)
	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		errCh <- err
	}

	// os.Getenv("SERVER_PORT")
	if err := svc.Init(
		micro.Server(httpsrv.NewServer(
			server.Name("file-service"),
			server.Version("1.0"),
			server.Address(":4545"),
			server.Context(ctx),
			server.Codec("application/json", jsoncodec.NewCodec()),
			server.WrapHandler(promwrapper.NewHandlerWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
			server.WrapHandler(idwrapper.NewServerHandlerWrapper()),
		)),
		micro.Client(httpcli.NewClient(
			client.ContentType("application/json"),
			client.Codec("application/json", jsoncodec.NewCodec()),
			client.Wrap(promwrapper.NewClientWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
			client.Wrap(idwrapper.NewClientWrapper()),
		)),
		micro.Store(s3),
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
		logger.Infof(ctx, "Not found, %v\n", r.URL)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("Not found. Path: %s, Method: %s", r.RequestURI, r.Method)))
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(ctx, "Method not allowed, %v\n", r.URL)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(fmt.Sprintf("Method not allowed, %s", r.Method)))
	})
	endpoints := pb.NewFileProcessingServiceEndpoints()

	// os.Getenv("DB_URL")
	db := <-initDB("postgres", "postgres://user:userpassword@localhost:5432/file_service_db?sslmode=disable", errCh)

	fr := repository.NewAWSFileRepository(db)

	srv := service.NewAWSProcessingService(jsoncodec.NewCodec(), fr, svc.Options().Store)

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
		stats.WithVersionDate("1.0", time.Now().String()),
	)

	healthServer := stats.NewServer(statsOpts...)
	go func() {
		logger.Fatal(ctx, healthServer.Serve(":9090"))
	}()

	if err := svc.Run(); err != nil {
		errCh <- err
	}
}
