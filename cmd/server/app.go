package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
	"turn_on_pc/internal/config"
	"turn_on_pc/internal/user"
	"turn_on_pc/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Infoln("create Router")
	router := httprouter.New()
	cfg := config.GetConfig()

	logger.Infoln("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Infoln("start application")

	var listenErr error
	var listener net.Listener
	if cfg.Listen.Type == "socket" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0])) // получение абсолютного пути где запущено приложение
		if err != nil {
			logger.Fatalln(err)
		}
		logger.Infoln("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Infoln("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)

	} else {
		logger.Infoln("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)) //todo сделать перезапись сокета
		logger.Infof("server is listening: http://%s:%s ", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatalln(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Fatalln(server.Serve(listener))

}
