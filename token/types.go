package token

// Credentials is the type used to pass username and password around.
type Credentials struct {
	User string
	Pass string
}

// credentialsDB maps the json db to struct.
type credentialsDB struct {
	userpassDB []struct {
		passwordDB string `json:"pass"`
		usernameDB string `json:"user"`
	} `json:"credentials"`
}
