package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/arieffian/mw-backend-test/internal/config"
	"github.com/arieffian/mw-backend-test/internal/connectors"
	helper "github.com/arieffian/mw-backend-test/pkg/helpers"

	log "github.com/sirupsen/logrus"
)

var (
	// Router instance
	Router *http.ServeMux

	// userLogger instance of logrus logger
	apiLogger = log.WithField("go", "API")

	// brandHandler http handler for brand routing
	brandHandler *BrandHandler

	// productHandler http handler for product routing
	productHandler *ProductHandler

	// transactionHandler http handler for transaction routing
	transactionHandler *TransactionHandler
)

func Start() {
	configureLogging()
	log.Infof("Starting User service")
	InitializeDB()
	InitializeRouter()

	var wait time.Duration

	address := fmt.Sprintf("%s:%s", config.Get("server.user.host"), config.Get("server.user.port"))
	log.Info("Server binding to ", address)

	srv := &http.Server{
		Addr:         address,
		Handler:      Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Infof("Shutting down")
	os.Exit(0)
}

func InitializeDB() {
	if config.Get("db.type") == "mysql" {
		log.Warnf("Using MYSQL")

		// init repo for token package
		BrandRepo = connectors.GetMySQLDBInstance()
		ProductRepo = connectors.GetMySQLDBInstance()
		TransactionRepo = connectors.GetMySQLDBInstance()
	} else {
		apiLogger.Fatal("unknown database type")
		panic(fmt.Sprintf("unknown database type %s. Correct your configuration 'db.type' or env-var 'AAA_DB_TYPE'. allowed values are INMEMORY or MYSQL", config.Get("db.type")))
	}

}

func configureLogging() {
	lLevel := config.Get("server.log.level")
	fmt.Println("Setting log level to ", lLevel)
	switch strings.ToUpper(lLevel) {
	default:
		fmt.Println("Unknown level [", lLevel, "]. Log level set to ERROR")
		log.SetLevel(log.ErrorLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	}

	lType := config.Get("log.type")
	fmt.Println("Setting log type to ", lType)
	if lType == "FILE" {
		helper.InitLogRotate()
		logFile := helper.GetFileLog()
		if logFile != nil {
			log.SetOutput(logFile)
		}
	}
}

func InitializeRouter() {
	log.Info("Initializing server")
	Router = http.NewServeMux()
	brandHandler = &BrandHandler{}
	productHandler = &ProductHandler{}
	transactionHandler = &TransactionHandler{}

	apiRoutes()
}

// initRouter will initialize router to execute API endpoints
func apiRoutes() {
	Router.HandleFunc("/brand", brandHandler.BrandHttpHandler)
	Router.HandleFunc("/product", productHandler.ProductHttpHandler)
	Router.HandleFunc("/product/brand", productHandler.ProductHttpHandler)
	Router.HandleFunc("/order", transactionHandler.TransactionHttpHandler)
}
