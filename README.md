# Annotate Registry Artifacts

![Main Branch Build Workflow Badge](https://github.com/johnsonshi/annotate-registry-artifacts/actions/workflows/build.yml/badge.svg)
![Latest Release Workflow Badge](https://github.com/johnsonshi/annotate-registry-artifacts/actions/workflows/release.yml/badge.svg)

Attach [OCI annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md) to _**existing**_ registry artifacts (such as [container images](https://www.docker.com/resources/what-container/) or [OCI artifacts](https://github.com/opencontainers/artifacts)) by attaching an annotations file using [ORAS Reference Types](https://oras.land/cli/6_reference_types/).

_**NOTE:**_ This only works for [registries supporting OCI Artifacts](https://oras.land/implementors/#registries-supporting-oci-artifacts), [ORAS Artifacts](https://github.com/oras-project/artifacts-spec), and [ORAS Artifact Reference Types](https://oras.land/cli/6_reference_types/).
This tool has been tested with [Azure Container Registry](https://azure.microsoft.com/en-us/services/container-registry/).

## Scenario

This tool is for you if you are a:

* Registry Owner, Maintainer, or Administrator
* Registry Artifacts – Container Image Builder, Maintainer, or Publisher

Registry owners and artifact publishers may wish to add custom [OCI annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md) existing images in the registry. Common annotation scenarios include:

* Annotation for image end-of-life date (EOL Date), which may or may not be known during image build time.
* Annotation for image deprecation (marking an image as deprecated).
* Annotation to note the date of a recent vulnerability scan.
* Annotation marking an image as an "official image", "golden image", "preferred image", or "premium image".
* Annotation for image compliance status, such as an image's compliance and certification to run in secure-cloud and government-cloud environments.

By design, directly modifying or adding [OCI annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md) to an _existing_ registry artifact is not possible after an artifact (such as a container image) has been built.
This is not possible as doing so would modify the hash digest of the existing registry artifact.

Additionally, various build tools for container images currently _do not_ support adding OCI Annotations _during_ Dockerfile image build.

This tool:

* creates a new annotation file containing the OCI Annotations you wish to add to an _existing_ registry artifact,
* pushes and stores the annotation file in the same registry and repository as the existing registry artifact,
* links the annotation file and the existing artifact using [ORAS Artifact References](https://oras.land/cli/6_reference_types/).

## Quick Start

### Install

To install, run the following commands.

```bash
curl -LO https://github.com/johnsonshi/annotate-registry-artifacts/releases/download/v0.0.1/annotation
chmod +x annotation
sudo mv annotation /usr/local/bin
```

### Attach

This command attaches a set of [OCI annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md) to an existing registry artifact (such as [container images](https://www.docker.com/resources/what-container/) or [OCI artifacts](https://github.com/opencontainers/artifacts)).

#### Attach – Usage

```bash
./bin/annotation attach \
  --username "$registry_username" \
  --password "$registry_password" \
  --registry "$registry_url" \
  --subject-repository "$repository_name" \
  --subject-tag-or-digest "$digest" \
  --annotation "org.opencontainers.image.source: https://www.github.com/user/repo/source" \
  --annotation "org.opencontainers.image.authors: EFGH Inc." \
  --annotation "org.opencontainers.image.vendor: ABCD Inc." \
  --annotation "org.opencontainers.image.licenses: ABCD Image License" \
  --annotation "com.example.image.custom.key1: val1" \
  --annotation "com.example.image.custom.key2: val2"
```

#### Attach – Result

![container-image-and-oras-artifact-manifest-with-oci-annotations-relationship](./docs/images/container-image-and-oras-artifact-manifest-with-oci-annotations-relationship.png)

## Additional Resources

For detailed explanations, please read the [detailed documentation page](./DETAILED_DOCS.md).
