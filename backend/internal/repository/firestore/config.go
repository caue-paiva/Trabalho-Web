package firestore

// FirebaseConfig holds Firebase-specific configuration loaded from YAML
type FirebaseConfig struct {
	ProjectID       string `yaml:"project_id"`
	CredentialsPath string `yaml:"credentials_path"`
}

// Collections holds the names of Firestore collections loaded from YAML
type Collections struct {
	Texts     string `yaml:"texts"`
	Images    string `yaml:"images"`
	Timelines string `yaml:"timelines"`
}
