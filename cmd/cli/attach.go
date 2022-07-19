/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	digest "github.com/opencontainers/go-digest"
	ocispecv1 "github.com/opencontainers/image-spec/specs-go/v1"
	artifactspecv1 "github.com/oras-project/artifacts-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type attachOpts struct {
	stdin              io.Reader
	stdout             io.Writer
	stderr             io.Writer
	username           string
	password           string
	registry           string
	subjectRepository  string
	subjectTagOrDigest string
	annotationSlice    []string
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

	f.StringVar(&opts.username, "username", "", "username to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("username")

	// TODO add support for --password-stdin (reading password from stdin) for more secure password input.
	f.StringVar(&opts.password, "password", "", "password to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("password")

	f.StringVar(&opts.registry, "registry", "", "hostname of the registry (example: myregistry.azurecr.io)")
	cobraCmd.MarkFlagRequired("registry")

	f.StringVar(&opts.subjectRepository, "subject-repository", "", "subject repository to attach annotations to")
	cobraCmd.MarkFlagRequired("subject-repository")

	f.StringVar(&opts.subjectTagOrDigest, "subject-tag-or-digest", "", "subject tag or digest (in the subject repository) to attach annotations to")
	cobraCmd.MarkFlagRequired("subject-tag-or-digest")

	f.StringArrayVar(&opts.annotationSlice, "annotation", []string{}, "annotation to attach to the subject reference artifact")
	cobraCmd.MarkFlagRequired("annotation")

	return cobraCmd
}

func (opts *attachOpts) run() error {
	ctx := context.Background()

	annotationsMap, err := getAnnotationsMap(opts.annotationSlice)
	if err != nil {
		return err
	}

	repo, err := opts.getAuthenticatedRemoteRepositoryClient()
	if err != nil {
		return err
	}

	subjectArtifactDescriptor, err := getArtifactDescriptorByRepositoryTagOrDigest(repo, opts.subjectTagOrDigest, ctx)
	if err != nil {
		return err
	}

	referenceManifest := artifactspecv1.Manifest{
		MediaType:    artifactspecv1.MediaTypeArtifactManifest,
		ArtifactType: "annotations/json",
		// Based on https://pkg.go.dev/oras.land/oras-go/v2/registry/remote@v2.0.0-rc.1#example-Repository.Push-ArtifactReferenceManifest
		// To push a ORAS reference manifest (with a subject artifact),
		// we must first download the subject artifact manifest from the registry,
		// obtain the subject artifact manifest's artifact descriptor,
		// and set the ORAS reference manifest subject field to the subject artifact descriptor.
		Subject:     subjectArtifactDescriptor,
		Annotations: annotationsMap,
	}
	referenceManifestContent, _ := json.Marshal(referenceManifest)
	referenceManifestDescriptor := ocispecv1.Descriptor{
		MediaType: artifactspecv1.MediaTypeArtifactManifest,
		Digest:    digest.FromBytes(referenceManifestContent),
		Size:      int64(len(referenceManifestContent)),
	}

	// Push the reference manifest descriptor and content
	err = repo.Push(ctx, referenceManifestDescriptor, bytes.NewReader(referenceManifestContent))
	if err != nil {
		return err
	}

	fmt.Println("Push finished")

	return nil
}

func (opts *attachOpts) getAuthenticatedRemoteRepositoryClient() (*remote.Repository, error) {
	// Create a client to the remote repository identified by a reference.
	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", opts.registry, opts.subjectRepository))
	if err != nil {
		return nil, err
	}

	// Set the repository auth credential client.
	repo.Client = &auth.Client{
		Credential: func(ctx context.Context, reg string) (auth.Credential, error) {
			return auth.Credential{
				Username: opts.username,
				Password: opts.password,
			}, nil
		},
	}

	return repo, nil
}

func getArtifactDescriptorByRepositoryTagOrDigest(repo *remote.Repository, tagOrDigest string, ctx context.Context) (*artifactspecv1.Descriptor, error) {
	ociDescriptor, rc, err := repo.FetchReference(ctx, tagOrDigest)
	if err != nil {
		return nil, err
	}

	defer rc.Close()
	pulled, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	// verify the fetched content
	if ociDescriptor.Size != int64(len(pulled)) || ociDescriptor.Digest != digest.FromBytes(pulled) {
		return nil, err
	}

	artifactDescriptor := artifactspecv1.Descriptor{
		MediaType: ociDescriptor.MediaType,
		Digest:    ociDescriptor.Digest,
		Size:      ociDescriptor.Size,
	}

	return &artifactDescriptor, nil
}

// getAnnotationsMap returns a map of annotations from a slice of annotation strings.
// strings in the slice should conform to the following format: "key: value".
func getAnnotationsMap(annotationSlice []string) (map[string]string, error) {
	re := regexp.MustCompile(`:\s*`)
	annotationsMap := make(map[string]string)
	for _, rawAnnotation := range annotationSlice {
		annotation := re.Split(rawAnnotation, 2)
		if len(annotation) != 2 {
			return nil, fmt.Errorf("invalid annotation: %s", rawAnnotation)
		}
		annotationsMap[annotation[0]] = annotation[1]
	}
	return annotationsMap, nil
}
