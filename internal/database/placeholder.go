package database

import "strconv"

// RewritePlaceholders converts `?` placeholders to PostgreSQL-style `$1, $2, ...`
// positional placeholders. SQLite natively supports `?` so this function is only
// needed for the PostgreSQL driver.
func RewritePlaceholders(query string) string {
	var result []byte
	n := 0
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(query); i++ {
		ch := query[i]

		// Track whether we are inside a quoted string to avoid rewriting
		// question marks that appear inside string literals.
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			result = append(result, ch)
			continue
		}
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			result = append(result, ch)
			continue
		}

		if ch == '?' && !inSingleQuote && !inDoubleQuote {
			n++
			result = append(result, '$')
			result = append(result, []byte(strconv.Itoa(n))...)
		} else {
			result = append(result, ch)
		}
	}

	return string(result)
}
