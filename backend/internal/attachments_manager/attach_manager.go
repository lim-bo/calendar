package attachmanager

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"log/slog"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrAlreadyExist = errors.New("attachment with such name already exist")
)

type AttachManager struct {
	cli     *minio.Client
	bucket  string
	host    string
	extHost string
	sqlcli  *pgxpool.Pool
}

type MinioCfg struct {
	Address    string
	User       string
	Pass       string
	BucketName string
	ExtAddress string
}

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	Options  map[string]string
}

func New(cfg *MinioCfg, sqlcfg *DBConfig) *AttachManager {
	optsStr := ""
	if len(sqlcfg.Options) != 0 {
		optsStr = "?"
		for k, v := range sqlcfg.Options {
			optsStr += k + "=" + v
		}
	}
	p, err := pgxpool.Connect(context.Background(), "postgresql://"+sqlcfg.User+":"+sqlcfg.Password+"@"+sqlcfg.Host+":"+sqlcfg.Port+"/"+sqlcfg.DBName+optsStr)
	if err != nil {
		log.Fatal(err)
	}
	client, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Pass, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &AttachManager{
		cli:     client,
		bucket:  cfg.BucketName,
		sqlcli:  p,
		extHost: cfg.ExtAddress,
		host:    cfg.Address,
	}
}

func (am *AttachManager) Load(file *models.FileLoad) error {
	_, err := am.cli.PutObject(context.Background(), am.bucket, file.Name, bytes.NewReader(file.Data), int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return errors.New("loading file error: " + err.Error())
	}
	return nil
}

func (am *AttachManager) LoadAttachment(eventID primitive.ObjectID, file *models.FileLoad) error {
	tx, err := am.sqlcli.Begin(context.Background())
	if err != nil {
		return errors.New("tx begining error: " + err.Error())
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `INSERT INTO attachments (event_id, name) VALUES ($1, $2);`, eventID.Hex(), file.Name)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ErrAlreadyExist
			}
		}
		return errors.New("inserting value error: " + err.Error())
	}
	err = am.Load(file)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return errors.New("commiting error: " + err.Error())
	}
	return nil
}

func (am *AttachManager) GetAttachments(eventID primitive.ObjectID) ([]*models.FileDownload, error) {
	rows, err := am.sqlcli.Query(context.Background(), `SELECT name FROM attachments WHERE event_id = $1;`, eventID.Hex())
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr == pgx.ErrNoRows {
			return nil, errors.New("there is no event with such id")
		}
		return nil, errors.New("getting files' names error: " + err.Error())
	}
	result := make([]*models.FileDownload, 0)
	for rows.Next() {
		var file models.FileDownload
		err = rows.Scan(&file.Name)
		if err != nil {
			return nil, errors.New("scanning value error: " + err.Error())
		}
		obj, err := am.cli.GetObject(context.Background(), am.bucket, file.Name, minio.GetObjectOptions{})
		if err != nil {
			return nil, errors.New("error getting object: " + err.Error())
		}
		raw, err := io.ReadAll(obj)
		if err != nil {
			return nil, errors.New("error reading object data: " + err.Error())
		}
		file.Data = raw
		slog.Debug("file", slog.String("content", string(file.Data)))
		result = append(result, &file)
		obj.Close()
	}
	return result, nil
}
