package main

type config struct {
	ApiConfig apiConfig `yaml:"api"`
	Jaeger jaegerConfig `yaml:"jaeger"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type jaegerConfig struct {
	URL string `yaml:"url"`
}