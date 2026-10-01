package model

type BackupConfig struct {
	Consoles []ConsoleConfig
}

func (config *BackupConfig) FindConsoleConfigByName(name string) *ConsoleConfig {
	for _, console := range config.Consoles {
		if console.Name == name {
			return &console
		}
	}
	return nil
}

