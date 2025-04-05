package models

type UserCredentials struct {
	Email string `json:"mail"`
	Pass  string `json:"pass"`
}

type UserCredentialsRegister struct {
	UserCredentials
	FirstName  string `json:"f_name"`
	SecondName string `json:"s_name"`
	ThirdName  string `json:"t_name"`
	Department string `json:"dep"`
	Position   string `json:"pos"`
}
