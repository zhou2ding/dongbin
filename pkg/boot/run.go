package boot

import (
	"blog/pkg/v"
	"blog/pkg/l"
	recover2 "blog/pkg/recover"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type AppRouter interface {
	RegisterRouters(router *mux.Router)
}

type Init func() error

type UnInit func()

type RunParam struct {
	AppName string
	Router  AppRouter
	Inits   []Init
	UnInits []UnInit
}

func ListenAndServe(p *RunParam) {
	if p.Inits != nil {
		for _, init := range p.Inits {
			if err := init(); err != nil {
				panic(err)
			}
		}
	}

	if p.UnInits != nil {
		defer func() {
			for _, unInit := range p.UnInits {
				unInit()
			}
		}()
	}

	port := v.GetViper().Sub(p.AppName).GetInt("port")
	if port <= 0 {
		return
	}
	router := newRouter()
	p.Router.RegisterRouters(router)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: recover2.NewRecoverHandler(true)(router),
	}
	defer server.Close()

	go func(s *http.Server) {
		l.GetLogger().Info("serving http", zap.Int("port", port))
		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				l.GetLogger().Info("http Sever is closed")
			} else {
				l.GetLogger().Fatal("can not start http server", zap.Error(err))
			}
		}
	}(&server)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		l.GetLogger().Info("receive interrupt signal")
	}

	l.GetLogger().Info("http service has stopped", zap.String("service", p.AppName))
}

func newRouter() *mux.Router {
	return mux.NewRouter().UseEncodedPath()
}
