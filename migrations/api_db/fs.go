package migrations

import "embed"

//go:embed 00001_init_schema.sql 00002_seed_roles.sql
var FS embed.FS
