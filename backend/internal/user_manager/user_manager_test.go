package usermanager_test

import (
	"testing"

	"github.com/google/uuid"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
)

var cfg usermanager.DBConfig

func TestRegister(t *testing.T) {
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
	cfg = usermanager.DBConfig{
		Host:     "localhost",
		Port:     "5435",
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	m.Run()
}

func TestLoginAndGetProfile(t *testing.T) {
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

func TestGetUUIDS(t *testing.T) {
	um := usermanager.New(cfg)
	uids, err := um.GetUUIDS([]string{"aaa", "mar@mail.ru"})
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range uids {
		t.Log(u)
	}
}

func TestGetMails(t *testing.T) {
	um := usermanager.New(cfg)
	mails, err := um.GetEmails([]uuid.UUID{uuid.MustParse("32bc8c40-d423-4607-88b6-42248f8d256f")})
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mails {
		t.Log(m)
	}
}
