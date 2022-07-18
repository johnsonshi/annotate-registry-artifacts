/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package main

import (
	"flag"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "annotation",
		Short: "Interact with the OCI annotations of a registry artifact",
	}

	cobraCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flags := cobraCmd.PersistentFlags()

	cobraCmd.AddCommand(
		newAttachCmd(stdin, stdout, stderr, args),
	)

	_ = flags.Parse(args)

	return cobraCmd
}

func execute() {
	rootCmd := newRootCmd(os.Stdin, os.Stdout, os.Stderr, os.Args[1:])
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
