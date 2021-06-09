module github.com/kucjac/cleango/database/es/esxsql

go 1.16

require (
	github.com/jmoiron/sqlx v1.3.4
	github.com/kucjac/cleango v0.0.14
	github.com/kucjac/cleango/database/xsql v0.0.14
	github.com/satori/go.uuid v1.2.0
)

replace github.com/kucjac/cleango => ../../../

replace github.com/kucjac/cleango/database/xsql => ../../xsql
