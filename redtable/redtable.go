package redtable

const COLUMN_FAMILY = "_values"
const STRING_VALUE_COLUMN = "value"
const EXPIRY_COLUMN = "exp"

// FQCN returns the fully-qualified column name
// of the given col. Format is family:col.
// Some real rocket science over here uh.
func FQCN(col string) string {
	return COLUMN_FAMILY + ":" + col
}
