package src

import (
	"slices"
)

type Scripts map[string]Script

type Script interface {
	GetScript(runner Runner) (string, bool)
}

type SimpleScript string

func (s SimpleScript) GetScript(_ Runner) (string, bool) {
	return string(s), true
}

type RunnerScript map[Runner]string

func (s RunnerScript) GetScript(runner Runner) (string, bool) {
	script, ok := s[runner]

	// if no exact match was found, try to find a Runner that should be compatible
	if !ok {
		// shell runners should all be able to run shell scripts
		if slices.Contains(ShellRunners, runner) {
			script, ok = s[Shell]
		}

		// powershell runners should be able to run windows powershell scripts
		if slices.Contains(PowershellRunners, runner) {
			script, ok = s[WindowsPowershell]
		}
	}

	return script, ok
}
