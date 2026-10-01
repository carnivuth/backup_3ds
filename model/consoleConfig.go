package model

type ConsoleConfig struct {
	Name string `yaml:"name"`
	User string `yaml:"user"`
	Password string `yaml:"password" `
	Port int `yaml:"port"`
	Dirs []string `yaml:"dirs"`
}
