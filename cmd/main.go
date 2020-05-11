package main

import (
	"flag"
	"fmt"
	"kkagitala/go-rest-api/middleware"
	"kkagitala/go-rest-api/pkg/oc"
	"kkagitala/go-rest-api/transport/pb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"

	ordersvc "kkagitala/go-rest-api/implementation"
	"kkagitala/go-rest-api/repository"
	"kkagitala/go-rest-api/service"
	"kkagitala/go-rest-api/transport"
	grpctransport "kkagitala/go-rest-api/transport/grpc"
	httptransport "kkagitala/go-rest-api/transport/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	_ "github.com/lib/pq"
)

const (
	port = ":50051"
)

func main() {
	// initialize our OpenCensus configuration and defer a clean-up
	//defer oc.Setup("order").Close()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowWarn())
		logger = log.With(logger,
			"svc", "order",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Start the monitoring task
	go middleware.NewMonitor(logger, 300)

	// Connecting to db
	/*type product struct {
		gorm.Model
		Code string
		Price uint
	}

	type Booking struct{
		Id      int    `json:”id”`
		User    string `json:”user”`
		Members int  `json:”members”`
	}*/

	var db *gorm.DB
	{
		var err error
		// Connect to the database
		db, err = gorm.Open("mysql", "root:kiaan@(localhost)/payments")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		defer db.Close()
	}
	// Migrate the schema
	//db.AutoMigrate(&product{})
	//db.AutoMigrate(&Booking{})

	// Create
	//var booking = Booking{Id: 1, User: "kiaan", Members: 10}
	//db.Create(&booking)

	// Create Order Service
	var svc service.Service
	{
		repository, err := repository.New(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = ordersvc.NewService(repository, logger)
		// Add service middleware here
		// Logging middleware
		svc = middleware.LoggingMiddleware(logger)(svc)
	}

	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)
		// Add endpoint level middlewares here
		// Trace server side endpoints with open census
		endpoints = transport.Endpoints{
			Create:  oc.ServerEndpoint("Create")(endpoints.Create),
			GetByID: oc.ServerEndpoint("GetByID")(endpoints.GetByID),
		}

	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		var (
			httpAddr = flag.String("http.addr", ":7070", "HTTP listen address")
		)
		flag.Parse()

		// Create Go kit endpoints for the Order Service
		// Then decorates with endpoint middlewares

		var h http.Handler
		{
			ocTracing := kitoc.HTTPServerTrace()
			serverOptions := []kithttp.ServerOption{ocTracing}
			h = httptransport.NewService(endpoints, serverOptions, logger)
		}

		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: h,
		}
		errs <- server.ListenAndServe()

	}()

	go func() {
		// set-up grpc transport
		var (
			ocTracing          = kitoc.GRPCServerTrace()
			serverOptions      = []kitgrpc.ServerOption{ocTracing}
			subscriptionServer = grpctransport.NewGRPCServer(endpoints, serverOptions, logger)
			grpcListener, _    = net.Listen("tcp", port)
			grpcServer         = grpc.NewServer()
		)

		var g group.Group
		{
			/*
				Add an actor (function) to the group.
				Each actor must be pre-emptable by an interrupt function.
				That is, if interrupt is invoked, execute should return.
				Also, it must be safe to call interrupt even after execute has returned.
				The first actor (function) to return interrupts all running actors.
				The error is passed to the interrupt functions, and is returned by Run.
			*/
			g.Add(func() error {
				logger.Log("transport", "gRPC", "addr", port)
				pb.RegisterSubscriptionServer(grpcServer, subscriptionServer)
				return grpcServer.Serve(grpcListener)
			}, func(error) {
				grpcListener.Close()
			})
		}

		{
			cancelInterrupt := make(chan struct{})
			g.Add(func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				select {
				case sig := <-c:
					return fmt.Errorf("received signal %s", sig)
				case <-cancelInterrupt:
					return nil
				}
			}, func(error) {
				close(cancelInterrupt)
			})
		}
		/*
			Run all actors (functions) concurrently. When the first actor returns,
			all others are interrupted. Run only returns when all actors have exited.
			Run returns the error returned by the first exiting actor
		*/
		level.Error(logger).Log("exit", g.Run())
	}()

	level.Error(logger).Log("exit", <-errs)
}
