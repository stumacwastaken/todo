package seed

import (
	"database/sql"
	_ "embed" //_ to keep it here and call the init function of embed package
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stumacwastaken/todo/log"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// import (
// 	"context"
// 	"io"
// 	"os"
// 	"os/signal"
// 	"syscall"
// )

// func Run(out, stderr io.Writer) error {
// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE)
// 	defer cancel()
// }

func init() {
	Cmd.PersistentFlags().StringVar(&DBConnString, "connstring", "", "mysql connection string to connect to")
	Cmd.PersistentFlags().StringVar(&LogLevel, "log-level", "info", "log level of the application. use error, warn, info, debug")
}

var (
	Cmd = &cobra.Command{
		Use:   "seed",
		Short: "seeds the todo database with todo items",
		Long: `seeds the todo database with a set of base todo items. This is only intended for developer 
			purposes at this time`,
		Run: migrate,
	}

	//go:embed sql/seed.sql
	seedDoc      string
	DBConnString string
	LogLevel     string
)

func migrate(cmd *cobra.Command, args []string) {

	//yes this includes the password. This is not a best practice.
	log.Default().Info("seeding database", zap.String("connection-string", DBConnString))
	//open connection
	db, err := sqlx.Connect("mysql", DBConnString)
	if err != nil {
		log.Default().Panic("failed to connect to database", zap.String("connection-string", DBConnString), zap.Error(err))
	}
	//open a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Default().Panic("failed to get transaction", zap.Error(err))
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}
			log.Default().Panic("transaction rollback failed. You should be concened", zap.Error(errTx))
			return
		}
	}()
	var res sql.Result
	//execute sql from the seed sql doc
	if res, err = tx.Exec(seedDoc); err != nil {
		log.Default().Error("failed to insert seed data", zap.Error(err))
		return
	}
	//commit it all
	if err := tx.Commit(); err != nil {
		log.Default().Error("failed to commit seed data", zap.Error(err))
		return
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		//I can't think of a single reasonable scenario why we'd be here.
		log.Default().Panic("You have absolutely no business being here", zap.Error(err))
	}
	log.Default().Info("seeded data", zap.Int("todos-inserted", int(affectedRows)))

}
