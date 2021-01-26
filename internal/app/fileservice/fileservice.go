package fileservice

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	fileconfig "github.com/unistack-org/micro-config-file/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	"github.com/vielendanke/file-service/configs"
	"github.com/vielendanke/file-service/internal/app/fileservice/handlers"
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
		config.NewConfig(
			config.Struct(cfg),
		),
		config.NewConfig(
			config.AllowFail(true),
			config.Struct(cfg),
			config.Codec(jsoncodec.NewCodec()),
			fileconfig.Path("configs/local.json"),
		),
	); err != nil {
		errCh <- err
	}
	options := append([]micro.Option{},
		micro.Server(httpsrv.NewServer()),
		micro.Client(httpcli.NewClient()),
		micro.Context(ctx),
		micro.Name("file-service"),
		micro.Version("1.0"),
	)
	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		errCh <- err
	}
	if err := svc.Init(
		micro.Server(httpsrv.NewServer(
			server.Name("file-service"),
			server.Version("1.0"),
			server.Address(":4545"),
			server.Context(ctx),
			server.Codec("application/json", jsoncodec.NewCodec()),
		)),
		micro.Client(httpcli.NewClient(
			client.ContentType("application/json"),
			client.Codec("application/json", jsoncodec.NewCodec()),
		)),
	); err != nil {
		errCh <- err
	}
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(ctx, "Not found, %v\n", r.URL)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(ctx, "Method not allowed, %v\n", r.URL)
	})
	endpoints := pb.NewFileProcessingServiceEndpoints()

	db := <-initDB("postgres", "postgres://user:userpassword@localhost:5432/file_service_db?sslmode=disable", errCh)

	srv := service.NewAWSProcessingService(jsoncodec.NewCodec(), db)

	handler := handlers.NewFileServiceHandler(srv, jsoncodec.NewCodec())

	if err := configs.ConfigureHandlerToEndpoints(router, handler, endpoints); err != nil {
		errCh <- err
	}
	if err := svc.Server().Handle(svc.Server().NewHandler(router)); err != nil {
		errCh <- err
	}
	if err := svc.Run(); err != nil {
		errCh <- err
	}
}
