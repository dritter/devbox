// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package boxcli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/devbox"
)

type shellCmdFlags struct {
	config   configFlags
	PrintEnv bool
}

func ShellCmd() *cobra.Command {
	flags := shellCmdFlags{}
	command := &cobra.Command{
		Use:   "shell -- [<cmd>]",
		Short: "Start a new shell or run a command with access to your packages",
		Long: "Start a new shell or run a command with access to your packages. \n" +
			"If invoked without `cmd`, devbox will start an interactive shell.\n" +
			"If invoked with a `cmd`, devbox will run the command in a shell and then exit.\n" +
			"In both cases, the shell will be started using the devbox.json found in the --config flag directory. " +
			"If --config isn't set, then devbox recursively searches the current directory and its parents.",
		Args:              validateShellArgs,
		PersistentPreRunE: nixShellPersistentPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShellCmd(cmd, args, flags)
		},
	}

	command.Flags().BoolVar(
		&flags.PrintEnv, "print-env", false, "Print script to setup shell environment")
	flags.config.register(command)
	return command
}

func runShellCmd(cmd *cobra.Command, args []string, flags shellCmdFlags) error {
	path, cmds, err := parseShellArgs(cmd, args, flags)
	if err != nil {
		return err
	}

	// Check the directory exists.
	box, err := devbox.Open(path, os.Stdout)
	if err != nil {
		return errors.WithStack(err)
	}

	if flags.PrintEnv {
		// return here to prevent opening a devbox shell
		return box.PrintShellEnv()
	}

	if devbox.IsDevboxShellEnabled() {
		return errors.New("You are already in an active devbox shell.\nRun 'exit' before calling devbox shell again. Shell inception is not supported.")
	}

	if len(cmds) > 0 {
		err = box.Exec(cmds...)
	} else {
		err = box.Shell()
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil
	}
	return err
}

func nixShellPersistentPreRunE(cmd *cobra.Command, args []string) error {
	_, err := exec.LookPath("nix-shell")
	if err != nil {
		return errors.New("could not find nix in your PATH\nInstall nix by following the instructions at https://nixos.org/download.html and make sure you've set up your PATH correctly")
	}
	return nil
}

func validateShellArgs(cmd *cobra.Command, args []string) error {
	lenAtDash := cmd.ArgsLenAtDash()
	if lenAtDash > 1 {
		return fmt.Errorf("accepts at most 1 directory, received %d", lenAtDash)
	}
	return nil
}

func parseShellArgs(cmd *cobra.Command, args []string, flags shellCmdFlags) (string, []string, error) {
	index := cmd.ArgsLenAtDash()
	if index < 0 {
		configPath, err := configPathFromUser(args, &flags.config)
		if err != nil {
			return "", nil, err
		}
		return configPath, []string{}, nil
	}

	path, err := configPathFromUser(args[:index], &flags.config)
	if err != nil {
		return "", nil, err
	}
	cmds := args[index:]

	return path, cmds, nil
}
