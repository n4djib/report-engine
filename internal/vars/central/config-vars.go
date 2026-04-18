package vars

type ConfigVars struct {
	AppName string `env:"APP_NAME" envDefault:"Report Engine Central"`
	AppUrl  string `env:"APP_URL" envDefault:"http://localhost"`
	AppPort int    `env:"APP_PORT" envDefault:"8080"`
	// AllowOrigins string `env:"ALLOW_ORIGINS" envDefault:"http://localhost:8080"`
	// ActivateUrl string `env:"ACTIVATION_URL,required"`
	// PasswordlessUrl    string `env:"PASSWORDLESS_URL,required"`
	// DBUrl              string `env:"DATABASE_URL" envDefault:"./database.db"`
	// Activati onExpMin   int    `env:"ACTIVATION_EXP_MINUTES" envDefault:"60"`
	// CookieExpMin       int    `env:"COOKIE_EXP_MINUTES" envDefault:"1440"` // 1 day
	// CookieSecure       bool   `env:"COOKIE_SECURE" envDefault:"true"`
	// PasswordlessExpMin int    `env:"PASSWORDLESS_EXP_MINUTES" envDefault:"15"`
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	HideBanner bool   `env:"HIDE_BANNER" envDefault:"false"`
	HidePort   bool   `env:"HIDE_PORT" envDefault:"false"`
	// BcryptSalt         int    `env:"BCRYPT_SALT" envDefault:"10"`
	// MinPasswordScore   int    `env:"MIN_PASSWORD_SCORE" envDefault:"0"`
	SMTP SMTPConfig
}

type SMTPConfig struct {
	Host string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	Port int    `env:"SMTP_PORT" envDefault:"587"`
	// User string `env:"SMTP_USER,required"`
	// Pass string `env:"SMTP_PASS,required"`
	// From string `env:"SMTP_FROM,required"`
}
