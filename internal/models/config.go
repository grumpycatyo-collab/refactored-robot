package models

type Config struct {
	Environment string `yaml:"environment"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
}
