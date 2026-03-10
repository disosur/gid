package internal

// Profile holds the git/SSH identity for a single GitHub account.
type Profile struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	SSHKey string `json:"sshKey"`
}

// Config is the top-level structure stored in ~/.gid/profiles.json.
type Config struct {
	Active   string             `json:"active"`
	Profiles map[string]Profile `json:"profiles"`
}
