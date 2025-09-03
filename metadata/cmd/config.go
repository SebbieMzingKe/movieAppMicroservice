package main

type serviceConfig struct {
	ApiConfig apiConfig `yaml:"api"`
	Jaeger jaegerConfig `yaml:"jaeger"`
}

type apiConfig struct {
	Port string `yaml:"port"`
}

type jaegerConfig struct {
	URL string `yaml:"url"`
}