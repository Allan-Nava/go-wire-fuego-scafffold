package app

import (
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
	
	
	"go.uber.org/zap"
)

type App struct {}

func (a *App) Routes(logger *zap.SugaredLogger) *fuego.Server {
    s := fuego.NewServer(
		fuego.WithAddr("0.0.0.0:8080"),
		fuego.WithCorsMiddleware(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		}).Handler),
		fuego.WithGlobalResponseTypes(400, "Bad ConsumerRequest _(validation or deserialization error)_", fuego.BadRequestError{}),
		fuego.WithGlobalResponseTypes(401, "Unauthorized _(authentication error)_", fuego.UnauthorizedError{}),
		//fuego.WithLogHandler(applog.NewLoggerAdapter(logger)),
	)

	fuego.Get(s, "/health", func(ctx fuego.ContextNoBody) (string, error) {
		return "OK", nil
	})

	return s
}