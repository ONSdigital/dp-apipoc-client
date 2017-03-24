package client

import "bytes"

func buildPath(fragments []string) string {
	var str bytes.Buffer

	for _, i := range fragments {
		str.WriteString(i)
	}

	return str.String()
}
