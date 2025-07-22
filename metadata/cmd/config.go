package main

type serviceConfig struct {
	ApiConfig apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port string `yaml:"port"`
}

type jaegerConfig struct {
	URL string `yaml:"url"`
}