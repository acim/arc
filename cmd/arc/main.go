package main

import (
	"fmt"
	"os"
	"time"

	"github.com/acim/arc/pkg/controller"
	"github.com/acim/arc/pkg/mail"
	"github.com/acim/arc/pkg/rest"
	"github.com/acim/arc/pkg/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/mailgun/mailgun-go/v4"
	"go.ectobit.com/act"
)

type dbConfig struct {
	Hostname string `def:"postgres"`
	Username string `def:"postgres"`
	Password string
	Name     string `def:"postgres"`
}

type config struct {
	ServiceName string `def:"arc"`
	ServerPort  int    `def:"3000"`
	MetricsPort int    `def:"3001"`
	Environment string `def:"dev"`
	JWT         struct {
		Secret                 string
		AuthTokenExpiration    time.Duration `env:"ARC_JWT_AUTH_TOKEN_EXP" def:"15m"`
		RefreshTokenExpiration time.Duration `env:"ARC_JWT_REFRESH_TOKEN_EXP" def:"168h"`
	}
	DB      dbConfig
	Mailgun struct {
		Domain    string
		APIKey    string
		Recipient string
	}
}

func main() { //nolint:funlen
	c := &config{} //nolint:exhaustivestruct

	if len(os.Args) < 2 { //nolint:gomnd
		usage()
	}

	cmd := os.Args[1]
	switch cmd {
	case "serve":
		serveCmd := act.New("serve", act.WithUsage("arc"))

		if err := serveCmd.Parse(c, os.Args[2:]); err != nil {
			exit("parse arguments", err)
		}

		logger, err := rest.NewLogger(c.Environment)
		if err != nil {
			exit("logger", err)
		}

		db, err := pgstore.NewDB(c.DB.Hostname, c.DB.Username, c.DB.Password, c.DB.Name)
		if err != nil {
			exit("connect to postgres", err)
		}

		users := pgstore.NewUsers(db, pgstore.UsersTableName("admin"))
		jwtAuth := jwtauth.New("HS256", []byte(c.JWT.Secret), nil)
		authController := controller.NewAuth(users, jwtAuth, logger)

		mailSender := mail.NewMailgun(mailgun.NewMailgun(c.Mailgun.Domain, c.Mailgun.APIKey))
		mailController := controller.NewMail(mailSender, c.Mailgun.Recipient, logger)

		router := rest.DefaultRouter(c.ServiceName, nil, logger)
		router.Post("/auth", authController.Login)
		router.Post("/mail", mailController.Send)

		router.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/auth", authController.User)
			r.Delete("/auth", authController.Logout)
		})

		// router.Get("/heavy", func(w http.ResponseWriter, r *http.Request) {
		// 	err := valve.Lever(r.Context()).Open()
		// 	if err != nil {
		// 		logger.Error("open valve lever", zap.Error(err))
		// 	}
		// 	defer valve.Lever(r.Context()).Close()

		// 	select {
		// 	case <-valve.Lever(r.Context()).Stop():
		// 		logger.Info("valve closed, finishing")
		// 	case <-time.After(2 * time.Second):
		// 		// Do some heave lifting
		// 		time.Sleep(2 * time.Second)
		// 	}

		// 	res := abmiddleware.ResponseFromContext(r.Context())
		// 	res.SetPayload("all done")
		// })

		app := rest.NewServer(c.ServiceName, c.ServerPort, c.MetricsPort, router, logger)
		app.Run()

	case "user":
		config := &struct{}{}
		userCmd := act.New("user", act.WithUsage("arc"))

		if err := userCmd.Parse(config, os.Args[2:]); err != nil {
			exit("parse arguments", err)
		}

	// 		email := add.Arg("email", "user's email address").Required().String()
	// password := add.Arg("password", "user's password").Required().String()
	// add.Action(func(c *kingpin.ParseContext) error {
	// 	if !govalidator.IsEmail(*email) {
	// 		return errors.New("email not valid")
	// 	}

	// 	if len(*password) < 8 {
	// 		return errors.New("password should contain minimum 8 characters")
	// 	}

	// 	user, err := model.NewUser(*email, *password)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return users.Insert(context.Background(), user)
	// })

	default:
		usage()
	}
}

func usage() {
	usage := `Usage of arc:
  arc <command>

  Available commands:
    serve	start rest server
    user	create new user`

	fmt.Println(usage) //nolint:forbidigo
	os.Exit(2)         //nolint:gomnd
}

func exit(message string, err error) {
	fmt.Printf("%s: %v\n", message, err) //nolint:forbidigo
	os.Exit(2)                           //nolint:gomnd
}
