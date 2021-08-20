module github.com/kucjac/cleango/database/es/esxsql

go 1.16

require (
	github.com/kucjac/cleango v0.0.24
	github.com/kucjac/cleango/database/xsql v0.0.24
	github.com/satori/go.uuid v1.2.0
)

replace github.com/kucjac/cleango => ../../../

replace github.com/kucjac/cleango/database/xsql => ../../xsql
