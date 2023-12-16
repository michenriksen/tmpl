package tmux

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// envArgs returns a slice of tmux command arguments to set the provided
// environment variables.
//
// The arguments are sorted by key to make the output deterministic.
func envArgs(envs ...map[string]string) []string {
	m := mergeMaps(envs...)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	args := make([]string, 0, len(m)*2)

	for _, k := range keys {
		eVal := fmt.Sprintf("%s=%s", k, os.ExpandEnv(m[k]))
		args = append(args, "-e", eVal)
	}

	return args
}

// mergeMaps merges the provided maps into a single map.
func mergeMaps(maps ...map[string]string) map[string]string {
	res := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}

	return res
}

// inTmux returns true if the application is running inside tmux.
func inTmux() bool {
	if os.Getenv("TERM_PROGRAM") == "tmux" {
		return true
	}

	if os.Getenv("TMUX") != "" {
		return true
	}

	if strings.Contains(os.Getenv("TERM"), "tmux") {
		return true
	}

	return false
}
