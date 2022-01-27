package main

type EmailConfig struct {
	EMAIL     string `mapstructure:"EMAIL"`
	PASSWORD  string `mapstructure:"PASSWORD"`
	SMTP_HOST string `mapstructure:"SMTP_HOST"`
	SMTP_PORT string `mapstructure:"SMTP_PORT"`
}
