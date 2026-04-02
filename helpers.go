package forge

import (
	"fmt"
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
	case time.Time:
		return strconv.Quote(x.Format(time.RFC3339))
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}
