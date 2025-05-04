package attachmanager_test

import (
	"testing"

	attachmanager "github.com/lim-bo/calendar/backend/internal/attachments_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
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
	am := attachmanager.New(cfg)
	file := models.FileLoad{
		Name: "file.txt",
		Data: []byte("file content"),
	}
	err := am.Load(&file)
	if err != nil {
		t.Error(err)
	}
}
