package configs

type ConfigClient interface {
	GetConfig(cfgName string) (any, error)
}

type configService struct {
}

func (s configService) GetConfig(cfgName string) (any, error) {
	return "", nil
}

func NewConfigService() ConfigClient {
	return configService{}
}
