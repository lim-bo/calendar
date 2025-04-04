package usermanager

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserManager struct {
	mu   *sync.RWMutex
	pool *pgxpool.Pool
}

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	Options  map[string]string
}

func New(cfg DBConfig) *UserManager {
	optsStr := ""
	if len(cfg.Options) != 0 {
		optsStr = "?"
		for k, v := range cfg.Options {
			optsStr += k + "=" + v
		}
	}
	p, err := pgxpool.Connect(context.Background(), "postgresql://"+cfg.User+":"+cfg.Password+"@"+cfg.Host+":"+cfg.Port+"/"+cfg.DBName+optsStr)
	if err != nil {
		log.Fatal(err)
	}
	return &UserManager{
		mu:   &sync.RWMutex{},
		pool: p,
	}
}
