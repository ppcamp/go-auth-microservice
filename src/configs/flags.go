package configs

import (
	"github.com/ppcamp/go-cli/env"
)

// Flags return all mappings to the config variables
func Flags() []env.Flag {
	return []env.Flag{
		//#region: Cache/Redis flags
		&env.BaseFlag[string]{
			EnvName: "CACHE_HOST",
			Value:   &CacheHost,
			Default: "localhost",
		},
		&env.BaseFlag[string]{
			EnvName: "CACHE_PORT",
			Value:   &CachePort,
			Default: "6379",
		},
		&env.BaseFlag[string]{
			Value:   &CachePassword,
			EnvName: "CACHE_PASSWORD",
		},
		&env.BaseFlag[int]{
			Value:   &CacheDb,
			Default: 0,
			EnvName: "CACHE_DATABASE",
		},
		//#endregion

		//#region: Database/Postgres flags
		&env.BaseFlag[string]{
			Value:   &DatabaseHost,
			Default: "localhost",
			EnvName: "DATABASE_HOST",
		},
		&env.BaseFlag[string]{
			Value:   &DatabasePort,
			Default: "5432",
			EnvName: "DATABASE_PORT",
		},
		&env.BaseFlag[string]{
			Value:   &DatabaseName,
			Default: "authentication",
			EnvName: "DATABASE_DBNAME",
		},
		&env.BaseFlag[string]{
			Value:   &DatabaseUser,
			Default: "authuser",
			EnvName: "DATABASE_USER",
		},
		&env.BaseFlag[string]{
			Value:   &DatabasePassword,
			Default: "somepassword",
			EnvName: "DATABASE_PASSWORD",
		},
		//#endregion

		//#region: App flags
		&env.BaseFlag[string]{
			Value:   &AppEnvironment,
			Default: "dev",
			EnvName: "APP_ENV",
		},
		&env.BaseFlag[string]{
			Value:   &AppEnvironment,
			Default: ":9000",
			EnvName: "APP_PORT",
		},
		&env.BaseFlag[string]{
			Value:   &AppEnvironment,
			Default: "3490be09e8904918997b073c460c834c",
			EnvName: "APP_ID",
		},
		//#endregion

		//#region: JWT/ Security flags
		&env.BaseFlag[string]{
			Value:   &JwtPublic,
			EnvName: "JWT_PUBLIC",
			Default: `ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBAFeQDDdKv96vAHOivPGDiEkSt02E8qFmGNz+aFv/hZ0fkP3QLeHxm7HCMZguNcCj3HFTi3HERj8jN0nHvTXbTM6TwEG6evRIzm8qZlzfl3CsBIRtiAFmdKKHHmO9sibd7gZZifU8+emZReFF1ZyYL0v5HuT8M2vs67J2vSsyTxjkyD4ww== ppcamp@DESKTOP-14OV55P
	`,
		},
		&env.BaseFlag[string]{
			Value:   &JwtPrivate,
			EnvName: "JWT_PRIVATE",
			Default: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAArAAAABNlY2RzYS
1zaGEyLW5pc3RwNTIxAAAACG5pc3RwNTIxAAAAhQQBXkAw3Sr/erwBzorzxg4hJErdNhPK
hZhjc/mhb/4WdH5D90C3h8ZuxwjGYLjXAo9xxU4txxEY/IzdJx70120zOk8BBunr0SM5vK
mZc35dwrASEbYgBZnSihx5jvbIm3e4GWYn1PPnpmUXhRdWcmC9L+R7k/DNr7Ouydr0rMk8
Y5Mg+MMAAAEYbnyzrm58s64AAAATZWNkc2Etc2hhMi1uaXN0cDUyMQAAAAhuaXN0cDUyMQ
AAAIUEAV5AMN0q/3q8Ac6K88YOISRK3TYTyoWYY3P5oW/+FnR+Q/dAt4fGbscIxmC41wKP
ccVOLccRGPyM3Sce9NdtMzpPAQbp69EjObypmXN+XcKwEhG2IAWZ0ooceY72yJt3uBlmJ9
Tz56ZlF4UXVnJgvS/ke5Pwza+zrsna9KzJPGOTIPjDAAAAQgFbCxbLXNkxBcdk+46SXOwr
x8tIUfjKNd+LZoiu7vFTk+V2L8jvaKlj3anxhcyrvSf28D8Jna1LZ5Ru+AaXgFLJCgAAAB
ZwcGNhbXBAREVTS1RPUC0xNE9WNTVQAQIDBA==
-----END OPENSSH PRIVATE KEY-----
`,
		},
		//#endregion
	}
}
