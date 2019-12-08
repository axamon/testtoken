package token

// Credentials is the type used to pass username and password around.
type Credentials struct {
	User string
	Pass string
}

// CredentialsDB maps the json db.
type CredentialsDB struct {
	Userpass []struct {
		Pass string `json:"pass"`
		User string `json:"user"`
	} `json:"credentials"`
}
