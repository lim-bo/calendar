package usermanager

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrRegistered   = errors.New("already registered")
	ErrUnregistered = errors.New("unregistered")
	ErrWrongPass    = errors.New("provided password is wrong")
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

func (um *UserManager) Register(creds *models.UserCredentialsRegister) error {

	hashpass, err := util.Hash(creds.Pass)
	if err != nil {
		return ErrInternal
	}
	uid := uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	um.mu.Lock()
	_, err = um.pool.Exec(ctx, `INSERT INTO profiles (uid, mail, password,  first_name, second_name, third_name, position, department) VALUES
($1, $2, $3, $4, $5, $6, $7, $8);`, uid, creds.Email, string(hashpass), creds.FirstName, creds.SecondName, creds.ThirdName, creds.Position, creds.Department)
	um.mu.Unlock()
	if err != nil {
		var pgerr pgx.PgError
		if errors.As(err, &pgerr) {
			return ErrRegistered
		}
		return err
	}
	return nil
}

func (um *UserManager) Login(creds *models.UserCredentials) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	um.mu.Lock()
	row := um.pool.QueryRow(ctx, `SELECT uid, password FROM profiles p WHERE p.mail = $1;`)
	um.mu.Unlock()
	var uidStr, hashpass string
	err := row.Scan(&uidStr, &hashpass)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.UUID{}, ErrUnregistered
		}
		return uuid.UUID{}, err
	}
	err = util.CheckPassword(creds.Pass, hashpass)
	if err != nil {
		return uuid.UUID{}, ErrWrongPass
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uid, nil
}
