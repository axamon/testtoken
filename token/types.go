package token

// Credentials is the type used to pass username and password around.
type Credentials struct {
	User string
	Pass string
}

// credentialsDB maps the json db to struct.
type credentialsDB struct {
	UserpassDB []struct {
		PasswordDB string `json:"pass"`
		UsernameDB string `json:"user"`
	} `json:"credentials"`
}
