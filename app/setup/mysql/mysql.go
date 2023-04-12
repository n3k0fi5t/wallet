package mysql

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

const (
	maxIdleConns = 10
	maxOpenConns = 128
	maxLifeTime  = 5 * time.Minute
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbName     = os.Getenv("DB_NAME")
	dbPort     = os.Getenv("DB_PORT")
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
)

var dbClient *sqlx.DB

func init() {
	db, err := sqlx.Open("mysql", getDSN())
	if err != nil {
		panic(err)
	}

	fmt.Printf("connect to %v\n", getDSN())

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(maxOpenConns)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(maxIdleConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(maxLifeTime)

	dbClient = db
}

func getDSN() string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPassword, dbHost, dbPort, dbName)
	return dsn
}

func GetMySQL() *sqlx.DB {
	return dbClient
}
