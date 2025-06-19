package main

type serviceConfig struct {
	ApiConfig apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port string `yaml:"port"`
}
