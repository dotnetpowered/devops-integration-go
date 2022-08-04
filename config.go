package main

type Config struct {
	//Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Zabbix        ZabbixConfig   `yaml:"zabbix"`
	OctopusDeploy OctopusConfig  `yaml:"octopus"`
	VmWare        VmWareConfig   `yaml:"vmware"`
	Database      DatabaseConfig `yaml:"database"`
}

type ZabbixConfig struct {
	Url      string `yaml:"url" env:"ZABBIX_URL"`
	User     string `yaml:"user" env:"ZABBIX_USER"`
	Password string `yaml:"password" env:"ZABBIX_PASSWORD"`
}

type OctopusConfig struct {
	Url    string `yaml:"url" env:"OCTOPUS_URL"`
	ApiKey string `yaml:"api-key" env:"OCTOPUS_APIKEY"`
}

type VmWareConfig struct {
	Url      string `yaml:"url" env:"VMWARE_URL"`
	User     string `yaml:"user" env:"VMWARE_USER"`
	Password string `yaml:"password" env:"VMWARE_PASSWORD"`
}

type DatabaseConfig struct {
	Dsn string `yaml:"dsn"`
}
