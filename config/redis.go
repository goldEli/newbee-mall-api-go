package config

type Redis struct {
	Host      string `mapstructure:"host" json:"host" yaml:"host"` // 服务器地址
	Port      string `mapstructure:"port" json:"port" yaml:"port"` // 端口
	DB        int    `mapstructure:"db" json:"db" yaml:"db"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
}

func (r *Redis) Addr() string {
	return r.Host + ":" + r.Port
}

// func (m *Mysql) Dsn() string {
// 	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
// }
