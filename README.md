# monobake

A Git tag-driven build helper for Docker monorepos.

## Motivation

While researching ways to build and push Docker images using GitHub Actions
and Git tags in a repository containing multiple applications and Dockerfiles,
I noticed something.

Nx, Bazel, and Turborepo are often discussed as monorepo solutions. However,
for someone who simply wants to build Docker images, these tools felt like
overkill in terms of setup effort and learning curve.

Writing shell scripts and jq directly in GitHub Actions is hard to maintain
and difficult to test locally.

I thought this middle layer was a gap waiting to be filled.

As a solution for this middle layer, I developed a tool that uses Docker Bake.
This tool generates parameters for Docker Bake — specifically, it resolves
build targets and versions from Git tags. The goals are:

* Eliminate logic from GitHub Actions YAML files
* Leave the integration with docker buildx bake to the caller
* Run the same command both locally and in GitHub Actions

## How It Works

monobake parses a Git tag like `backend/v1.0.0` and extracts:

- Target: `backend`
- Version: `1.0.0`

Pass the target to `docker buildx bake` to build only that target.
Use `--set` to apply the version as an image tag.

## Usage

```bash
# Resolve build target from tag
$ monobake -tag refs/tags/backend/v1.0.0
backend 1.0.0

# Combine with Bake
read TARGET VERSION <<< $(monobake -tag refs/tags/backend/v1.0.0)
[ -n "$TARGET" ] && docker buildx bake \
  --set="${TARGET}.tags=${REGISTRY}/${TARGET}:${VERSION}" \
  "$TARGET" --push
```

## Options

```
-tag string     Git tag to parse
-file string    Path to Bake file (default: docker-bake.json)
-version        Show version
```

## Bake File Format

```json
{
  "target": {
    "backend": {
      "context": "apps/backend",
      "dockerfile": "Dockerfile"
    },
    "frontend": {
      "context": "apps/frontend",
      "dockerfile": "Dockerfile"
    }
  }
}
```

## Git Tag Format

```
<target>/v<major>.<minor>.<patch>
```

Examples: `backend/v1.0.0`, `frontend/v2.1.0-beta.1`

## GitHub Actions

See [monobake-demo](https://github.com/zinrai/monobake-demo) for a working example.

## Exit Codes

| Code | Meaning             |
|------|---------------------|
| 0    | Success             |
| 1    | Configuration error |

Returns exit code 0 with no output if the tag format is invalid.

## License

This project is licensed under the [MIT License](./LICENSE).
