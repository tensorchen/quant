package config

type Config struct {
	Tquant     Tquant     `yaml:"tquant"`
	LongBridge LongBridge `yaml:"long_bridge"`
	Discord    Discord    `yaml:"discord"`
}

type Tquant struct {
	Port  string `yaml:"port"`
	Token string `yaml:"token"`
}

type LongBridge struct {
	AppKey      string `yaml:"app_key"`
	AppSecret   string `yaml:"app_secret"`
	AccessToken string `yaml:"access_token"`
}

type Discord struct {
	ID    string `yaml:"id"`
	Token string `yaml:"token"`
}
