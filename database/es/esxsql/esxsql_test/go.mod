module github.com/kucjac/cleango/database/es/esxsql/esxsql_test

go 1.16

require (
	github.com/kucjac/cleango v0.0.14
	github.com/kucjac/cleango/database/es/esxsql v0.0.14
	github.com/kucjac/cleango/database/xpq v0.0.14
	github.com/kucjac/cleango/database/xsql v0.0.14
	github.com/lib/pq v1.10.2
)

replace (
	github.com/kucjac/cleango => ./../../../../
	github.com/kucjac/cleango/database/es/esxsql => ./../
	github.com/kucjac/cleango/database/xpq => ./../../../xpq
	github.com/kucjac/cleango/database/xsql => ./../../../xsql
)
