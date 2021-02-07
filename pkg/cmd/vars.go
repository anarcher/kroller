package cmd

import "strings"

type stringSlice []string

func (ss *stringSlice) String() string {
	if len(*ss) <= 0 {
		return "..."
	}

	return strings.Join(*ss, ", ")
}

func (ss *stringSlice) Set(value string) error {
	*ss = append(*ss, value)

	return nil
}
