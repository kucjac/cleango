module github.com/kucjac/cleango/database/es/esxsql/esxsql_test

go 1.16

require (
	github.com/kucjac/cleango v0.0.20
	github.com/kucjac/cleango/database/es/esxsql v0.0.20
	github.com/kucjac/cleango/database/xpq v0.0.20
	github.com/kucjac/cleango/database/xsql v0.0.20
	github.com/lib/pq v1.10.2
)

replace github.com/kucjac/cleango => ../../../../

replace github.com/kucjac/cleango/database/es/esxsql => ./..

replace github.com/kucjac/cleango/database/xpq => ../../../xpq

replace github.com/kucjac/cleango/database/xsql => ../../../xsql
