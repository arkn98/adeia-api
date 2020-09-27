/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import "adeia/pkg/util/constants"

// envOverrides holds all environment value keys for overriding the config.
var envOverrides = map[string]string{
	"server.jwt_secret": constants.EnvServerJWTSecretKey,

	"mailer.username": constants.EnvMailerUsernameKey,
	"mailer.password": constants.EnvMailerPasswordKey,

	"database.dbname":   constants.EnvDBNameKey,
	"database.user":     constants.EnvDBUserKey,
	"database.password": constants.EnvDBPasswordKey,
	"database.host":     constants.EnvDBHostKey,
	"database.port":     constants.EnvDBPortKey,

	"cache.host": constants.EnvCacheHostKey,
	"cache.port": constants.EnvCachePortKey,
}

// Config represents the overall configuration.
type Config struct {
	CacheConfig  `mapstructure:"cache"`
	DBConfig     `mapstructure:"database"`
	LoggerConfig `mapstructure:"logger"`
	MailerConfig `mapstructure:"mailer"`
	ServerConfig `mapstructure:"server"`
}

// CacheConfig represents the config for the cache.
type CacheConfig struct {
	Network  string `mapstructure:"network"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	ConnSize int    `mapstructure:"connsize"`
}

// DBConfig represents the config for the database.
type DBConfig struct {
	Driver      string `mapstructure:"driver"`
	DBName      string `mapstructure:"dbname"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	SSLMode     string `mapstructure:"sslmode"`
	SSLCert     string `mapstructure:"sslcert,omitempty"`
	SSLKey      string `mapstructure:"sslkey,omitempty"`
	SSLRootCert string `mapstructure:"sslrootcert,omitempty"`
}

// LoggerConfig represents the config for the logger.
type LoggerConfig struct {
	Level string   `mapstructure:"level"`
	Paths []string `mapstructure:"paths"`
}

// MailerConfig represents the config for the mailer.
type MailerConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
}

// ServerConfig represents the config for the server.
type ServerConfig struct {
	Host            string `mapstructure:"host,omitempty"`
	Port            int    `mapstructure:"port"`
	RateLimitRate   int    `mapstructure:"ratelimit_rate"`
	RateLimitWindow int    `mapstructure:"ratelimit_window"`
	JWTSecret       string `mapstructure:"jwt_secret"`
}
