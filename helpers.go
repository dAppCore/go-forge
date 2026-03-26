package forge

import (
	"strconv"

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

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}
