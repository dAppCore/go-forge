package forge

import (
	"fmt"
	// Note: AX-6 intrinsic — query escaping is centralised here because the pinned core module has no core.URLEncode helper.
	"net/url"
	"sort"
	// Note: AX-6 intrinsic — the pinned core module has no core.Itoa/core.FormatInt/core.Atoi/core.FormatBool/core.Quote helpers.
	"strconv"
	"strings"
	"time"

	core "dappco.re/go/core"
)

func trimTrailingSlashes(s string) string {
	for core.HasSuffix(s, "/") {
		s = core.TrimSuffix(s, "/")
	}
	return s
}

func int64String(v int64) string {
	return strconv.FormatInt(v, 10)
}

func intString(v int) string {
	return strconv.Itoa(v)
}

func boolString(v bool) string {
	return strconv.FormatBool(v)
}

func parseInt(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}

func pathParams(values ...string) Params {
	params := make(Params, len(values)/2)
	for i := 0; i+1 < len(values); i += 2 {
		params[values[i]] = values[i+1]
	}
	return params
}

func optionString(typeName string, fields ...any) string {
	var b strings.Builder
	b.WriteString(typeName)
	b.WriteString("{")

	wroteField := false
	for i := 0; i+1 < len(fields); i += 2 {
		name, _ := fields[i].(string)
		value := fields[i+1]
		if isZeroOptionValue(value) {
			continue
		}
		if wroteField {
			b.WriteString(", ")
		}
		wroteField = true
		b.WriteString(name)
		b.WriteString("=")
		b.WriteString(formatOptionValue(value))
	}

	b.WriteString("}")
	return b.String()
}

func isZeroOptionValue(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case bool:
		return !x
	case int:
		return x == 0
	case int64:
		return x == 0
	case []string:
		return len(x) == 0
	case *time.Time:
		return x == nil
	case *bool:
		return x == nil
	case time.Time:
		return x.IsZero()
	default:
		return false
	}
}

func formatOptionValue(v any) string {
	switch x := v.(type) {
	case string:
		return strconv.Quote(x)
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case []string:
		return fmt.Sprintf("%#v", x)
	case *time.Time:
		if x == nil {
			return "<nil>"
		}
		return strconv.Quote(x.Format(time.RFC3339))
	case *bool:
		if x == nil {
			return "<nil>"
		}
		return strconv.FormatBool(*x)
	case time.Time:
		return strconv.Quote(x.Format(time.RFC3339))
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func serviceString(typeName, fieldName string, value any) string {
	return typeName + "{" + fieldName + "=" + fmt.Sprint(value) + "}"
}

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

type queryBuilder struct {
	values map[string][]string
}

func newQueryBuilder() *queryBuilder {
	return &queryBuilder{values: make(map[string][]string)}
}

func (q *queryBuilder) Set(key, value string) {
	q.values[key] = []string{value}
}

func (q *queryBuilder) Add(key, value string) {
	q.values[key] = append(q.values[key], value)
}

func (q *queryBuilder) Encode() string {
	if q == nil || len(q.values) == 0 {
		return ""
	}

	keys := make([]string, 0, len(q.values))
	for key := range q.values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, key := range keys {
		escapedKey := url.QueryEscape(key)
		for _, value := range q.values[key] {
			if b.Len() > 0 {
				b.WriteByte('&')
			}
			b.WriteString(escapedKey)
			b.WriteByte('=')
			b.WriteString(url.QueryEscape(value))
		}
	}
	return b.String()
}

func encodeQuery(build func(*queryBuilder)) string {
	query := newQueryBuilder()
	build(query)
	return query.Encode()
}

func appendQuery(path string, build func(*queryBuilder)) string {
	query := encodeQuery(build)
	if query == "" {
		return path
	}
	if strings.Contains(path, "?") {
		return path + "&" + query
	}
	return path + "?" + query
}
