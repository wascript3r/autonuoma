package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	// Logger
	_loggerUcase "github.com/wascript3r/cryptopay/pkg/logger/usecase"

	// User
	_userHandler "github.com/wascript3r/autonuoma/pkg/user/delivery/http"
	_userWsHandler "github.com/wascript3r/autonuoma/pkg/user/delivery/ws"
	_userPwHasher "github.com/wascript3r/autonuoma/pkg/user/pwhasher"
	_userRepo "github.com/wascript3r/autonuoma/pkg/user/repository"
	_userUcase "github.com/wascript3r/autonuoma/pkg/user/usecase"
	_userValidator "github.com/wascript3r/autonuoma/pkg/user/validator"

	// Session
	_sessionMid "github.com/wascript3r/autonuoma/pkg/session/delivery/http/middleware"
	_sessionWsMid "github.com/wascript3r/autonuoma/pkg/session/delivery/ws/middleware"
	_sessionGen "github.com/wascript3r/autonuoma/pkg/session/generator"
	_sessionRepo "github.com/wascript3r/autonuoma/pkg/session/repository"
	_sessionUcase "github.com/wascript3r/autonuoma/pkg/session/usecase"

	// Message
	_messageWsHandler "github.com/wascript3r/autonuoma/pkg/message/delivery/ws"
	_messageEventBus "github.com/wascript3r/autonuoma/pkg/message/eventbus"
	_messageRepo "github.com/wascript3r/autonuoma/pkg/message/repository"
	_messageUcase "github.com/wascript3r/autonuoma/pkg/message/usecase"
	_messageValidator "github.com/wascript3r/autonuoma/pkg/message/validator"

	// Review
	_reviewHandler "github.com/wascript3r/autonuoma/pkg/review/delivery/http"
	_reviewRepo "github.com/wascript3r/autonuoma/pkg/review/repository"
	_reviewUcase "github.com/wascript3r/autonuoma/pkg/review/usecase"
	_reviewValidator "github.com/wascript3r/autonuoma/pkg/review/validator"

	// Ticket
	_ticketWsHandler "github.com/wascript3r/autonuoma/pkg/ticket/delivery/ws"
	_ticketWsMid "github.com/wascript3r/autonuoma/pkg/ticket/delivery/ws/middleware"
	_ticketEventBus "github.com/wascript3r/autonuoma/pkg/ticket/eventbus"
	_ticketRepo "github.com/wascript3r/autonuoma/pkg/ticket/repository"
	_ticketUcase "github.com/wascript3r/autonuoma/pkg/ticket/usecase"
	_ticketValidator "github.com/wascript3r/autonuoma/pkg/ticket/validator"

	// License
	_licenseHandler "github.com/wascript3r/autonuoma/pkg/license/delivery/http"
	_licenseRepo "github.com/wascript3r/autonuoma/pkg/license/repository"
	_licenseUcase "github.com/wascript3r/autonuoma/pkg/license/usecase"
	_licenseValidator "github.com/wascript3r/autonuoma/pkg/license/validator"

	// Trip
	_tripHandler "github.com/wascript3r/autonuoma/pkg/trip/delivery/http"
	_tripRepo "github.com/wascript3r/autonuoma/pkg/trip/repository"
	_tripUcase "github.com/wascript3r/autonuoma/pkg/trip/usecase"

	// Reservation
	_reservationHandler "github.com/wascript3r/autonuoma/pkg/reservation/delivery/http"
	_reservationRepo "github.com/wascript3r/autonuoma/pkg/reservation/repository"
	_reservationUcase "github.com/wascript3r/autonuoma/pkg/reservation/usecase"

	// FAQ
	_faqHandler "github.com/wascript3r/autonuoma/pkg/faq/delivery/http"
	_faqRepo "github.com/wascript3r/autonuoma/pkg/faq/repository"
	_faqUcase "github.com/wascript3r/autonuoma/pkg/faq/usecase"

	// Room
	_roomRepo "github.com/wascript3r/autonuoma/pkg/room/repository"
	_roomUcase "github.com/wascript3r/autonuoma/pkg/room/usecase"

	// CORS
	_corsMid "github.com/wascript3r/autonuoma/pkg/cors/delivery/http/middleware"

	// Cars
	_carsHandler "github.com/wascript3r/autonuoma/pkg/cars/delivery/http"
	_carsRepo "github.com/wascript3r/autonuoma/pkg/cars/repository"
	_carsUcase "github.com/wascript3r/autonuoma/pkg/cars/usecase"
	_carsValidator "github.com/wascript3r/autonuoma/pkg/cars/validator"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/gocipher/aes"
	"github.com/wascript3r/gopool"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/eventbus"
	wsMiddleware "github.com/wascript3r/gows/middleware"
	_socketPool "github.com/wascript3r/gows/pool"
	"github.com/wascript3r/gows/router"
	"github.com/wascript3r/httputil/middleware"
)

const (
	// Database
	DatabaseDriver = "postgres"

	// Pool
	PoolSize   = 32
	WsPoolSize = 512

	// BTC
	AppLoggerPrefix = "[APP]"

	// WebSockets
	WSNetwork = "tcp"
)

var (
	WorkDir string
	Cfg     *Config

	flagLicensesDir = flag.String("img", "public/licenses/", "license images directory path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	// Get working directory
	WorkDir, err = os.Getwd()
	if err != nil {
		fatalError(err)
	}

	// Parse config file
	cfgPath, err := getConfigPath()
	if err != nil {
		fatalError(err)
	}

	Cfg, err = parseConfig(filepath.Join(WorkDir, cfgPath))
	if err != nil {
		fatalError(err)
	}
}

func fatalError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if *flagLicensesDir == "" {
		log.Fatal("Missing license images directory path (-img public/licenses/)")
	}

	// Logging
	logFlags := 0
	if Cfg.Log.ShowTimestamp {
		logFlags = log.Ltime
	}
	logger := _loggerUcase.New(
		AppLoggerPrefix,
		log.New(os.Stdout, "", logFlags),
	)

	// Startup message
	logger.Info("... Starting app ...")

	// Goroutine pool
	pool := gopool.New(PoolSize, 0, 0)
	wsPool := gopool.New(WsPoolSize, 0, 0)

	// Database connection
	dbConn, err := openDatabase(DatabaseDriver, Cfg.Database.Postgres.DSN)
	if err != nil {
		fatalError(err)
	}

	// AES cipher
	cipher, err := aes.NewCipher(Cfg.Cipher.AES.Key)
	if err != nil {
		fatalError(err)
	}

	// Session
	sessionRepo := _sessionRepo.NewPgRepo(dbConn)
	sessionGen := _sessionGen.New()
	sessionUcase := _sessionUcase.New(
		sessionRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		sessionGen,
		cipher,
		_sessionUcase.SessionLifetime(Cfg.Auth.Session.SessionLifetime.Duration),
	)

	// User
	userRepo := _userRepo.NewPgRepo(dbConn)
	userPwHasher := _userPwHasher.New(Cfg.Auth.PasswordCost)
	userValidator := _userValidator.New(userRepo)
	userUcase := _userUcase.New(
		userRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		sessionUcase,
		userPwHasher,
		userValidator,
	)

	// Message, Ticket
	messageRepo := _messageRepo.NewPgRepo(dbConn)
	ticketRepo := _ticketRepo.NewPgRepo(dbConn)

	// Message
	messageEventBus := _messageEventBus.New(pool, logger)
	messageValidator := _messageValidator.New()
	messageUcase := _messageUcase.New(
		messageRepo,
		ticketRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		messageEventBus,
		messageValidator,
	)

	// Review
	reviewRepo := _reviewRepo.NewPgRepo(dbConn)
	reviewValidator := _reviewValidator.New()
	reviewUcase := _reviewUcase.New(
		reviewRepo,
		ticketRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		reviewValidator,
	)

	// Ticket
	ticketEventBus := _ticketEventBus.New(pool, logger)
	ticketValidator := _ticketValidator.New(messageValidator)
	ticketUcase := _ticketUcase.New(
		ticketRepo,
		messageRepo,
		reviewRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		ticketEventBus,
		ticketValidator,
	)

	// License
	licenseRepo := _licenseRepo.NewPgRepo(dbConn)
	licenseValidator := _licenseValidator.New()
	licenseUcase := _licenseUcase.New(
		licenseRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		licenseValidator,

		"license",
		*flagLicensesDir,
	)

	// Trip
	tripRepo := _tripRepo.NewPgRepo(dbConn)
	tripUsecase := _tripUcase.New(tripRepo)

	// Reservation
	reservationRepo := _reservationRepo.NewPgRepo(dbConn)
	reservationUcase := _reservationUcase.New(reservationRepo)

	// FAQ
	faqRepo := _faqRepo.NewPgRepo(dbConn)
	faqUcase := _faqUcase.New(
		faqRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,
	)

	// Cars
	carsRepo := _carsRepo.NewPgRepo(dbConn)
	carsValidator := _carsValidator.New()
	carsUcase := _carsUcase.New(
		carsRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,
		carsValidator,
	)

	// Room
	roomRepo := _roomRepo.NewMemoryRepo()
	roomUcase := _roomUcase.New(roomRepo)

	// WS
	wsEventBus := eventbus.NewWsEventBus(wsPool, logger)

	wsRouter := router.New(wsEventBus)

	sessionWsMid := _sessionWsMid.NewWSMiddleware(
		_sessionWsMid.DefaultSessionKey,
		sessionUcase,
	)

	authWsStack := wsMiddleware.New()
	authWsStack.Use(sessionWsMid.Authenticated)

	notAuthWsStack := wsMiddleware.New()
	notAuthWsStack.Use(sessionWsMid.NotAuthenticated)

	clientWsStack := wsMiddleware.New()
	clientWsStack.Use(sessionWsMid.HasRole(domain.ClientRole))

	agentWsStack := wsMiddleware.New()
	agentWsStack.Use(sessionWsMid.HasRole(domain.AgentRole))

	adminWsStack := wsMiddleware.New()
	adminWsStack.Use(sessionWsMid.HasRole(domain.AdminRole))

	// App context
	ctx, cancel := context.WithCancel(context.Background())

	socketPool, err := _socketPool.NewPool(ctx, wsPool, logger, wsEventBus)
	if err != nil {
		fatalError(err)
	}

	ticketWsMid := _ticketWsMid.NewWSMiddleware(socketPool)

	_userWsHandler.NewWSHandler(
		wsRouter,
		notAuthWsStack,

		userUcase,
		sessionUcase,
		sessionWsMid,
		roomUcase,

		socketPool,
	)
	_messageWsHandler.NewWSHandler(
		wsRouter,
		clientWsStack,
		agentWsStack,

		messageUcase,
		messageEventBus,
		sessionUcase,
		ticketWsMid,

		socketPool,
	)
	_ticketWsHandler.NewWSHandler(
		wsRouter,
		clientWsStack,
		agentWsStack,

		ticketUcase,
		ticketEventBus,
		ticketWsMid,
		sessionUcase,
		roomUcase,

		socketPool,
	)

	wsListener, err := net.Listen(WSNetwork, ":"+Cfg.WebSocket.Port)
	if err != nil {
		fatalError(err)
	}

	wsServer, err := gows.NewServer(
		wsPool,
		wsListener,
		wsEventBus,

		gows.ConnIdleTime(Cfg.WebSocket.ConnIdleTime.Duration),
	)
	if err != nil {
		fatalError(err)
	}

	// Graceful shutdown
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := wsServer.Start(ctx); err != nil {
		fatalError(err)
	}

	// HTTP server
	httpRouter := httprouter.New()
	httpRouter.MethodNotAllowed = MethodNotAllowedHnd
	httpRouter.NotFound = NotFoundHnd
	httpRouter.ServeFiles("/licenses/*filepath", http.Dir(*flagLicensesDir))

	if Cfg.HTTP.EnablePprof {
		// pprof
		httpRouter.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)
	}

	sessionMid := _sessionMid.NewHTTPMiddleware(
		Cfg.Auth.Session.CookieName,
		Cfg.Auth.Session.CookieLifetime.Duration,
		Cfg.Auth.Session.SecureCookie,

		sessionUcase,
	)

	authStack := middleware.NewCtx()
	authStack.Use(sessionMid.Authenticated)

	notAuthStack := middleware.New()
	notAuthStack.Use(sessionMid.NotAuthenticated)

	clientStack := middleware.NewCtx()
	clientStack.Use(sessionMid.HasRole(domain.ClientRole))

	agentStack := middleware.NewCtx()
	agentStack.Use(sessionMid.HasRole(domain.AgentRole))

	_userHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		authStack,
		notAuthStack,

		userUcase,
		sessionUcase,
		sessionMid,
	)
	_reviewHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		clientStack,

		reviewUcase,
		sessionUcase,
	)
	_licenseHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		agentStack,
		clientStack,

		licenseUcase,
		sessionUcase,
	)

	_tripHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		clientStack,

		tripUsecase,
	)

	_reservationHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		sessionUcase,
		sessionMid,
		clientStack,
		reservationUcase,
	)

	_faqHandler.NewHTTPHandler(httpRouter, faqUcase)

	_carsHandler.NewHTTPHandler(httpRouter, carsUcase)

	httpServer := &http.Server{
		Addr: ":" + Cfg.HTTP.Port,
		Handler: _corsMid.NewHTTPMiddleware(
			Cfg.HTTP.CORS.Origin,
		).EnableCors(httpRouter),
	}

	// websocket.html file
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/websocket.html")
	})

	// Graceful shutdown
	gracefulShutdown := func() {
		cancel()
		pool.Terminate()
		wsPool.Terminate()

		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("Cannot shutdown HTTP server: %s", err)
		}

		logger.Info("... Exited ...")
		os.Exit(0)
	}

	go func() {
		<-stopSig
		logger.Info("... Received stop signal ...")
		gracefulShutdown()
	}()

	if err := httpServer.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			fmt.Println(err)
			gracefulShutdown()
		}
	}

	var c chan struct{}
	<-c
}
