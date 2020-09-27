/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package constants

const (
	// MaxReqBodySize (in bytes; default: 1MiB)
	MaxReqBodySize = 1048576
	// APIVersion represents the current major version of the API. It is used as URL prefix.
	APIVersion = "v1"

	// EmployeeIDLength represents the length of the generated employee IDs.
	EmployeeIDLength = 6
	// EmployeeIDChars represents the list of possible characters that can occur in an employee ID.
	EmployeeIDChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// ==========
	// Keys of env variables to override the config
	// ==========

	// EnvPrefix is used as the prefix for all env variables related to adeia.
	EnvPrefix = "ADEIA"

	// EnvServerJWTSecretKey is the env key for server's jwt secret.
	EnvServerJWTSecretKey = EnvPrefix + "_SERVER_JWT_SECRET"

	// Mailer keys

	// EnvMailerUsernameKey is the env key for mailer username.
	EnvMailerUsernameKey = EnvPrefix + "_MAILER_USERNAME"
	// EnvMailerPasswordKey is the env key for mailer password.
	EnvMailerPasswordKey = EnvPrefix + "_MAILER_PASSWORD"

	// Database (Postgres) keys

	// EnvConfPathKey is the env key for confPath.
	EnvConfPathKey = EnvPrefix + "_CONF_PATH"
	// EnvDBNameKey is the env key for database name.
	EnvDBNameKey = EnvPrefix + "_DB_NAME"
	// EnvDBUserKey is the env key for database user.
	EnvDBUserKey = EnvPrefix + "_DB_USER"
	// EnvDBPasswordKey is the env key for database password.
	EnvDBPasswordKey = EnvPrefix + "_DB_PASSWORD"
	// EnvDBHostKey is the env key for database host.
	EnvDBHostKey = EnvPrefix + "_DB_HOST"
	// EnvDBPortKey is the env key for database port.
	EnvDBPortKey = EnvPrefix + "_DB_PORT"

	// Cache (redis) keys

	// EnvCacheHostKey is the env key for redis host.
	EnvCacheHostKey = EnvPrefix + "_CACHE_HOST"
	// EnvCachePortKey is the env key for redis port.
	EnvCachePortKey = EnvPrefix + "_CACHE_PORT"
)
