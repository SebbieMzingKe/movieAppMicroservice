package main

type config struct {
	ApiConfig  apiConfig        `yaml:"api"`
	Jaeger     jaegerConfig     `yaml:"jaeger"`
	Prometheus prometheusConfig `yaml:"prometheus"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type jaegerConfig struct {
	URL string `yaml:"url"`
}

type prometheusConfig struct {
	URL string `yaml:"url"`
}
