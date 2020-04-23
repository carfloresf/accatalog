package main // Package main implements the go-service command.
import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hellerox/AcCatalog/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/hellerox/AcCatalog/pkg/handlers/rest"
	"github.com/hellerox/AcCatalog/pkg/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = "user=acadmin dbname=accat sslmode=disable"
	}

	// register healthcheck handler
	healthcheckHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	http.Handle("/healthcheck", healthcheckHandler)

	// initialize service
	service := service.AcCatalogService{Storage: storage.NewStorage(connectionString)}

	// initialize handlers
	handler := rest.MakeHTTPHandlers(&service)
	http.Handle("/", handler)

	// run http service & wait
	runService(cast.ToInt(port))
}

func runService(port int) {
	// init server
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", port),
	}

	errc := make(chan error, 2)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logrus.WithField("port", port).Info("listening")
		errc <- server.ListenAndServe()
	}()

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("accat initialized...")

	logrus.WithFields(logrus.Fields{
		"reason": <-errc,
	}).Info("terminated")
}
