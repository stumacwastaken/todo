package rest

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/stumacwastaken/todo/log"
	"github.com/stumacwastaken/todo/rest"
	"github.com/stumacwastaken/todo/stores/database"
	"github.com/stumacwastaken/todo/stores/tododb"
	"github.com/stumacwastaken/todo/todoitem"
	"github.com/stumacwastaken/todo/tracing"
	"go.uber.org/zap"
)

var (
	Cmd = &cobra.Command{
		Use:   "server",
		Short: "starts the todo restful server",
		Long:  ``,
		Run:   server,
	}
	Address  string
	Port     string
	LogLevel string
	DBConfig database.Config
)

func init() {
	Cmd.PersistentFlags().StringVar(&Address, "addr", "0.0.0.0", "address for the rest server to listen on")
	Cmd.PersistentFlags().StringVar(&Port, "port", "9000", "port for the server to listen on")
	Cmd.PersistentFlags().StringVar(&LogLevel, "log-level", "info", "log level of the application. use error, warn, info, debug")
	Cmd.PersistentFlags().StringVar(&DBConfig.Host, "dbhost", "localhost:3306", "mysql host and port")
	Cmd.PersistentFlags().StringVar(&DBConfig.User, "dbuser", "", "mysql user")
	Cmd.PersistentFlags().StringVar(&DBConfig.Password, "dbpass", "", "mysql password")
	Cmd.PersistentFlags().StringVar(&DBConfig.Name, "dbname", "todo", "database name")
}

func server(cmd *cobra.Command, args []string) {
	//make logger
	newLogger, err := log.New(LogLevel)
	if err != nil {
		panic(err) //we don't want to be running without logs. Fullstop
	}
	log.SetDefault(newLogger)

	//create database connection
	db, err := database.Open(DBConfig)
	if err != nil {
		//throw a panic if you can't connect to the database on startup. Likely a config issue.
		log.Default().Panic("failed to connect to database. Are your configs correct?", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		log.Default().Panic("failed to ping database.....", zap.Error(err))
	}
	srv := rest.NewServer(Address, Port)

	tdh := rest.NewTodoHandlers(todoitem.NewCore(tododb.NewStore(db)))
	//give a default base path for this server of api for now. It's entirely possible we can do this in networking though with k8s
	//basically, be ready to refactor and rip out
	tdh.RegisterTodoEndpoints(srv.Router, "/api")

	//register tracing
	tp := tracing.InitTracingProvider("todo")
	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()

	err = srv.Start(ctx)
	if err != nil {
		log.Default().Error("error starting restful server. Shutting down", zap.Error(err))
	}

}
