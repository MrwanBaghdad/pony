package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const injectHelp = `Assign a secret to a subshell`

func (cmd *injectCommand) Name() string              { return "inject" }
func (cmd *injectCommand) Args() string              { return "[KEYS] --sh [ SHELL COMMAND ]" }
func (cmd *injectCommand) ShortHelp() string         { return createHelp }
func (cmd *injectCommand) LongHelp() string          { return createHelp }
func (cmd *injectCommand) Hidden() bool              { return false }
func (cmd *injectCommand) Register(fs *flag.FlagSet) {}

type injectCommand struct {
	commndArgs []string
	secrets    []secret
}

type secret struct {
	key   string
	value string
}

func (cmd *injectCommand) Run(ctx context.Context, args []string) error {
	argsDoubleDashIdx := findDoubleDashIndex(args)
	numSecrets := argsDoubleDashIdx

	if numSecrets <= 0 {
		return errors.New("Must pass a secret")
	}

	secretKeys := args[:argsDoubleDashIdx]

	rootCmd := args[argsDoubleDashIdx+1]
	cmdArgs := args[argsDoubleDashIdx+2:]

	subCmd := exec.Command(rootCmd, cmdArgs...)

	subCmd.Stdout = os.Stdout
	subCmd.Stderr = os.Stderr
	subCmd.Env = os.Environ()

	for _, key := range secretKeys {
		if val, ok := s.Secrets[key]; ok {
			// Split the key for dotted keys and take last part to be injected
			splitted := strings.Split(key, ".")
			lastSuffix := splitted[len(splitted)-1]
			subCmd.Env = append(
				subCmd.Env,
				fmt.Sprintf("%s=%s", lastSuffix, val),
			)
		} else {
			return fmt.Errorf("secret for key %s does not exist", key)
		}
	}

	subCmd.Run()

	return nil
}

func findDoubleDashIndex(args []string) int {
	for i, char := range args {
		if char == "--" {
			return i
		}
	}
	return -1
}
