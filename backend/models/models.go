package models

type UserCredentials struct {
	Email string `json:"mail"`
	Pass  string `json:"pass"`
}

type UserCredentialsRegister struct {
	UserCredentials `json:",inline"`
	FirstName       string `json:"f_name"`
	SecondName      string `json:"s_name"`
	ThirdName       string `json:"t_name,omitempty"`
	Department      string `json:"dep"`
	Position        string `json:"pos"`
}
