package host

type Config struct {
	Hosts []Host `yaml:"hosts"`
	Do    string `yaml:"do"`
	Undo  string `yaml:"undo"`
}

type Host struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
