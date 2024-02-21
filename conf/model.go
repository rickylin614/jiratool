package conf

type BaseConfig struct {
	AdminPath     string `yaml:"admin_path"`
	AdminFileName string `yaml:"admin_filename"`
	WebPath       string `yaml:"web_path"`
	WebFileName   string `yaml:"source_filename"`
}
