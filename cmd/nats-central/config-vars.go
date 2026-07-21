package main

type ConfigVars struct {
	AppName      string `env:"APP_NAME" envDefault:"Central-NATS"`
	AppEnv       string `env:"APP_ENV" envDefault:"development"`
	NatsNkeySeed string `env:"NATS_NKEY_SEED,required"`
	NatsURL      string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
}
