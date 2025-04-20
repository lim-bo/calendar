package usermanager_test

import (
	"testing"

	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
)

func TestRegister(t *testing.T) {

	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	um := usermanager.New(cfg)
	testCreds := &models.UserCredentialsRegister{
		UserCredentials: models.UserCredentials{
			Email: "aaa",
			Pass:  "bbb",
		},
		FirstName:  "Ivan",
		SecondName: "Ivanov",
		ThirdName:  "Ivanovich",
		Department: "Some Department",
		Position:   "cleaner",
	}
	err := um.Register(testCreds)
	if err != nil {
		if err == usermanager.ErrRegistered {
			t.Log("user registered")
		} else {
			t.Error(err)
		}
	}
}

func TestLogin(t *testing.T) {
	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	um := usermanager.New(cfg)
	testCreds := &models.UserCredentials{
		Email: "testmail@gmail.com",
		Pass:  "secretPassword",
	}
	uid, err := um.Login(testCreds)
	if err != nil {
		switch err {
		case usermanager.ErrUnregistered:
			t.Log("user unregistered")
		case usermanager.ErrWrongPass:
			t.Log("wrong pass")
		case usermanager.ErrInternal:
			fallthrough
		default:
			t.Fatal(err)
		}
	}
	t.Log("uuid: ", uid)
}

func TestLoginAndUpdate(t *testing.T) {
	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}

	um := usermanager.New(cfg)
	uid, err := um.Login(&models.UserCredentials{
		Email: "testmail@gmail.com",
		Pass:  "secretPassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	testCreds := &models.UserCredentialsRegister{
		FirstName:  "NeIvan",
		SecondName: "NeIvanov",
		ThirdName:  "NeIvanovich",
		Department: "Some Department",
		Position:   "cleaner",
	}
	err = um.UpdateUser(testCreds, uid)
	if err != nil {
		t.Error(err)
	}
}

func TestLoginAndChangePassword(t *testing.T) {
	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	um := usermanager.New(cfg)
	uid, err := um.Login(&models.UserCredentials{
		Email: "testmail@gmail.com",
		Pass:  "secretPassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = um.ChangePassword("nonsecretPassword", uid)
	if err != nil {
		t.Error(err)
	}

}

func TestMain(m *testing.M) {
	util.LoadConfig()
	m.Run()
}

func TestLoginAndGetProfile(t *testing.T) {
	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	um := usermanager.New(cfg)
	uid, err := um.Login(&models.UserCredentials{
		Email: "testmail@gmail.com",
		Pass:  "secretPassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	info, err := um.GetProfileInfo(uid)
	if err != nil {
		t.Error(err)
	}
	t.Log(*info)
}
