package main

import "strings"

func parseEnv(envs []string) map[string]string {
	env := make(map[string]string, len(envs))

	for _, s := range envs {
		parts := strings.SplitN(s, "=", 2)
		env[parts[0]] = parts[1]
	}

	return env
}
