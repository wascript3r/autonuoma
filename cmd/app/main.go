package main

import (
	"context"
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
	_messageRepo "github.com/wascript3r/autonuoma/pkg/message/repository"
	_messageUcase "github.com/wascript3r/autonuoma/pkg/message/usecase"
	_messageValidator "github.com/wascript3r/autonuoma/pkg/message/validator"

	// Ticket
	_ticketWsHandler "github.com/wascript3r/autonuoma/pkg/ticket/delivery/ws"
	_ticketWsMid "github.com/wascript3r/autonuoma/pkg/ticket/delivery/ws/middleware"
	_ticketEventBus "github.com/wascript3r/autonuoma/pkg/ticket/eventbus"
	_ticketRepo "github.com/wascript3r/autonuoma/pkg/ticket/repository"
	_ticketUcase "github.com/wascript3r/autonuoma/pkg/ticket/usecase"
	_ticketValidator "github.com/wascript3r/autonuoma/pkg/ticket/validator"

	// Room
	_roomRepo "github.com/wascript3r/autonuoma/pkg/room/repository"
	_roomUcase "github.com/wascript3r/autonuoma/pkg/room/usecase"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/gocipher/aes"
	"github.com/wascript3r/gopool"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/eventbus"
	wsMiddleware "github.com/wascript3r/gows/middleware"
	_socketPool "github.com/wascript3r/gows/pool"
	"github.com/wascript3r/gows/router"
	"github.com/wascript3r/httputil"
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
	messageValidator := _messageValidator.New()
	messageUcase := _messageUcase.New(
		messageRepo,
		ticketRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		messageValidator,
	)

	// Ticket
	ticketEventBus := _ticketEventBus.New(pool, logger)
	ticketValidator := _ticketValidator.New(messageValidator)
	ticketUcase := _ticketUcase.New(
		ticketRepo,
		messageRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		ticketEventBus,
		ticketValidator,
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
	wsLog := func(next router.Handler) router.Handler {
		return func(ctx context.Context, s *gows.Socket, r *router.Request) {
			defer next(ctx, s, r)

			ss, err := sessionUcase.LoadCtx(ctx)
			if err != nil {
				log.Println("[WS] cannot get user ID")
				return
			}
			log.Println("[WS] user ID:", ss.UserID)
		}
	}
	authWsStack.Use(wsLog)

	notAuthWsStack := wsMiddleware.New()
	notAuthWsStack.Use(sessionWsMid.NotAuthenticated)

	clientWsStack := wsMiddleware.New()
	clientWsStack.Use(sessionWsMid.HasRole(domain.UserRole))

	supportWsStack := wsMiddleware.New()
	supportWsStack.Use(sessionWsMid.HasRole(domain.SupportRole))

	adminWsStack := wsMiddleware.New()
	adminWsStack.Use(sessionWsMid.HasRole(domain.AdminRole))
	adminWsStack.Use(wsLog)

	ticketWsMid := _ticketWsMid.NewWSMiddleware(
		_ticketWsMid.DefaultRoomKey,
		ticketUcase,
	)

	// App context
	ctx, cancel := context.WithCancel(context.Background())

	socketPool, err := _socketPool.NewPool(ctx, wsPool, logger, wsEventBus)
	if err != nil {
		fatalError(err)
	}

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

		messageUcase,
		sessionUcase,
		socketPool,
	)
	_ticketWsHandler.NewWSHandler(
		wsRouter,
		clientWsStack,
		supportWsStack,

		ticketUcase,
		sessionUcase,
		ticketWsMid,
		ticketEventBus,
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
	authStack.Use(func(next httputil.HandleCtx) httputil.HandleCtx {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			defer next(ctx, w, r, p)

			s, err := sessionUcase.LoadCtx(ctx)
			if err != nil {
				log.Println("cannot get user ID")
				return
			}
			log.Println("user ID:", s.UserID)
		}
	})

	notAuthStack := middleware.New()
	notAuthStack.Use(sessionMid.NotAuthenticated)

	_userHandler.NewHTTPHandler(
		context.Background(),

		httpRouter,
		authStack,
		notAuthStack,

		userUcase,
		sessionUcase,
		sessionMid,
	)

	httpServer := &http.Server{
		Addr:    ":" + Cfg.HTTP.Port,
		Handler: httpRouter,
	}

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
