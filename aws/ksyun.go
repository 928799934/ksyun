package aws

type Auth struct {
	AccessKey string
	SecretKey string
	TokenKey  string
}

type Region struct {
	Name       string
	S3Endpoint string
	DisableSSL bool
}
