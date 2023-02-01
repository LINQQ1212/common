package core

type System struct {
	Env       string `mapstructure:"env" json:"env" yaml:"env"`                      // 环境值
	Addr      int    `mapstructure:"addr" json:"addr" yaml:"addr"`                   // 端口值
	AdminAddr int    `mapstructure:"admin-addr" json:"admin-addr" yaml:"admin-addr"` // 端口值
	UserName  string `mapstructure:"username" json:"username" yaml:"username"`
	PassWord  string `mapstructure:"password" json:"password" yaml:"password"`
	MainDir   string `mapstructure:"main-dir" json:"main-dir" yaml:"main-dir"`
	ImageGrcp string `mapstructure:"image-grcp" json:"image-grcp" yaml:"image-grcp"`
	TGToken   string `mapstructure:"tg-token" json:"tg-token" yaml:"tg-token"`
	TGChatId  string `mapstructure:"tg-chat-id" json:"tg-chat-id" yaml:"tg-chat-id"`
	Info      string `mapstructure:"info" json:"info" yaml:"info"`
}
