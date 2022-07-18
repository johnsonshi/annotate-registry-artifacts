/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type attachOpts struct {
	stdin           io.Reader
	stdout          io.Writer
	stderr          io.Writer
	username        string
	password        string
	registry        string
	subject         string
	annotationSlice []string
}

func newAttachCmd(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) *cobra.Command {
	opts := &attachOpts{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	cobraCmd := &cobra.Command{
		Use:     "attach",
		Short:   "TODO",
		Example: `TODO`,
		RunE: func(_ *cobra.Command, args []string) error {
			return opts.run()
		},
	}

	f := cobraCmd.Flags()

	f.StringVarP(&opts.username, "username", "u", "", "username to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("username")

	// TODO add support for --password-stdin (reading password from stdin) for more secure password input.
	f.StringVarP(&opts.password, "password", "p", "", "password to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("password")

	f.StringVarP(&opts.registry, "registry", "r", "", "v2 API url of the registry (such as https://myregistry.azurecr.io/v2/")
	cobraCmd.MarkFlagRequired("registry")

	f.StringVarP(&opts.subject, "subject", "s", "", "subject image reference of the config annotation")
	cobraCmd.MarkFlagRequired("subject")

	f.StringArrayVarP(&opts.annotationSlice, "annotation", "a", []string{}, "annotation to add to the generated manifest's config")
	cobraCmd.MarkFlagRequired("annotation")

	return cobraCmd
}

func (opts *attachOpts) run() error {
	client := &auth.Client{
		Credential: func(ctx context.Context, reg string) (auth.Credential, error) {
			return auth.Credential{
				Username: opts.username,
				Password: opts.password,
			}, nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, opts.registry, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)

	return nil
}
