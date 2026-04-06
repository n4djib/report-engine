package main

type ConfigVars struct {
	AppName    string `env:"APP_NAME" envDefault:"Report Engine Remote"`
	AppUrl     string `env:"APP_URL" envDefault:"http://localhost"`
	AppPort    int    `env:"APP_PORT" envDefault:"8081"`
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	HideBanner bool   `env:"HIDE_BANNER" envDefault:"false"`
}
