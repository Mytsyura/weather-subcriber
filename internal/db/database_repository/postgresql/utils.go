package postgresql

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // must be imported
)

func GetDialect() goqu.DialectWrapper {
	return goqu.Dialect("postgres")
}
