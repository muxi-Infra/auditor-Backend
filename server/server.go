package server

import (
	"context"
	"net/http"
	"time"

	"github.com/muxi-Infra/auditor-Backend/controller"
	"github.com/muxi-Infra/auditor-Backend/middleware"
	"github.com/muxi-Infra/auditor-Backend/server/router"
)

type Server struct {
	Srv   *http.Server
	close func()
}

func NewServer(
	OAuth *controller.AuthController,
	User *controller.UserController,
	Item *controller.ItemController,
	Tube *controller.TubeController,
	Project *controller.ProjectController,
	LLm *controller.LLMController,
	Remove *controller.RemoveController,
	AuthMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CorsMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
) *Server {
	return &Server{
		Srv: &http.Server{
			Handler: router.NewRouter(OAuth, User, Project, Item, LLm, Remove, Tube,
				AuthMiddleware, corsMiddleware, loggerMiddleware),
		},
		close: func() {
			LLm.Close()
		},
	}
}

func (srv *Server) Run(addr string) error {
	srv.Srv.Addr = addr

	if err := srv.Srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (srv *Server) Close() {
	srv.close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Srv.Shutdown(ctx)
}
