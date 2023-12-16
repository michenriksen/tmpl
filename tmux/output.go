package tmux

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	newline = []byte("\n")
	comma   = []byte(",")
	colon   = []byte(":")
)

// outputRecord represents a line of output from a tmux command, that follows
// a specific format for parsing it as key-value pairs.
type outputRecord map[string]string

// String returns a string representation of the output record.
//
// NOTE: The order of the fields is not guaranteed to be consistent.
func (r outputRecord) String() string {
	res := make([]string, 0, len(r))

	for k, v := range r {
		res = append(res, fmt.Sprintf("%s:%s", k, v))
	}

	return strings.Join(res, ",")
}

// outputFormat returns a tmux command output format to be used with the -F
// flag.
//
// The vars are expected to be valid tmux variable names (e.g. session_id).
//
// The output format is a comma-separated list of key-value pairs, where the
// key is the variable name and the value is the variable placeholder:
//
//	"session_id:#{session_id},session_name:#{session_name}"
func outputFormat(vars ...string) string {
	res := make([]string, 0, len(vars))

	for _, v := range vars {
		res = append(res, fmt.Sprintf("%s:%s", v, outputFormatVar(v)))
	}

	return strings.Join(res, ",")
}

// parseOutput parses the tmux command output into a slice of output records.
//
// The output is expected to follow the format created by the [outputFormat]
// function.
func parseOutput(output []byte) ([]outputRecord, error) {
	output = bytes.TrimSpace(output)
	if len(output) == 0 {
		return []outputRecord{}, nil
	}

	lines := bytes.Split(output, newline)
	res := make([]outputRecord, 0, len(lines))

	for _, line := range lines {
		line = bytes.TrimSpace(line)

		if len(lines) == 0 {
			continue
		}

		record := make(outputRecord)
		keyvals := bytes.Split(line, comma)

		for _, kv := range keyvals {
			key, val, ok := bytes.Cut(kv, colon)
			if !ok {
				return nil, fmt.Errorf("invalid key-value pair in command output: %s", kv)
			}

			if _, ok := record[string(key)]; ok {
				return nil, fmt.Errorf("duplicate key in command output: %s", key)
			}

			record[string(key)] = string(val)
		}

		res = append(res, record)
	}

	return res, nil
}

func outputFormatVar(name string) string {
	return fmt.Sprintf("#{%s}", name)
}
