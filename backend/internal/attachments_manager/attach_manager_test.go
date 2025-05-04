package attachmanager_test

import (
	"testing"

	attachmanager "github.com/lim-bo/calendar/backend/internal/attachments_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {
	util.LoadConfig()
	m.Run()
}

func TestLoadFile(t *testing.T) {
	cfg := attachmanager.MinioCfg{
		Address:    viper.GetString("minio_addr"),
		User:       viper.GetString("minio_user"),
		Pass:       viper.GetString("minio_pass"),
		BucketName: viper.GetString("minio_bucket"),
	}
	sqlcfg := attachmanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	am := attachmanager.New(&cfg, &sqlcfg)
	file := models.FileLoad{
		Name: "file.txt",
		Data: []byte("file content"),
	}
	err := am.Load(&file)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadAttachment(t *testing.T) {
	cfg := attachmanager.MinioCfg{
		Address:    viper.GetString("minio_addr"),
		User:       viper.GetString("minio_user"),
		Pass:       viper.GetString("minio_pass"),
		BucketName: viper.GetString("minio_bucket"),
	}
	sqlcfg := attachmanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	am := attachmanager.New(&cfg, &sqlcfg)
	file := models.FileLoad{
		Name: "file2.txt",
		Data: []byte("secret"),
	}
	eventID, err := primitive.ObjectIDFromHex("6814a8011117e998968fcc97")
	if err != nil {
		t.Fatal(err)
	}
	err = am.LoadAttachment(eventID, &file)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAttachments(t *testing.T) {
	cfg := attachmanager.MinioCfg{
		Address:    viper.GetString("minio_addr"),
		User:       viper.GetString("minio_user"),
		Pass:       viper.GetString("minio_pass"),
		BucketName: viper.GetString("minio_bucket"),
	}
	sqlcfg := attachmanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	am := attachmanager.New(&cfg, &sqlcfg)
	eventID, err := primitive.ObjectIDFromHex("6814a8011117e998968fcc97")
	if err != nil {
		t.Fatal(err)
	}
	files, err := am.GetAttachments(eventID)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		t.Logf("filename: %s link: %s", f.Name, f.Link)
	}
}
