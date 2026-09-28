package engine

type BackupConfig struct {
	Consoles []ConsoleConfig
}

type ConsoleConfig struct {
	Name string `yaml:"name"`
	Port int `yaml:"port"`
	Dirs []string `yaml:"dirs"`
}
