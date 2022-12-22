package config

import "github.com/LINQQ1212/common/config/core"

type Server struct {
	JWT        core.JWT          `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap        core.Zap          `mapstructure:"zap" json:"zap" yaml:"zap"`
	System     core.System       `mapstructure:"system" json:"system" yaml:"system"`
	SeaFS      core.SeaFS        `mapstructure:"sea-fs" json:"sea-fs" yaml:"sea-fs"`
	Local      core.Local        `mapstructure:"local" json:"local" yaml:"local"`
	Domains    map[string]string `mapstructure:"domains" json:"domains" yaml:"domains"`
	Statistics map[string]string `mapstructure:"statistics" json:"statistics" yaml:"statistics"`
}
