# k8s-yaml-splitter
It takes a combined kubernetes yaml config and splits it into multiple files in a folder of your choosing

## Installation

### Download Binary
Download the latest binary for your platform from the [releases page](https://github.com/ohauer/k8s-yaml-splitter/releases):

```console
# Linux
curl -L -o k8s-yaml-splitter https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/k8s-yaml-splitter-linux-amd64
chmod +x k8s-yaml-splitter

# macOS
curl -L -o k8s-yaml-splitter https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/k8s-yaml-splitter-darwin-amd64
chmod +x k8s-yaml-splitter

# FreeBSD
curl -L -o k8s-yaml-splitter https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/k8s-yaml-splitter-freebsd-amd64
chmod +x k8s-yaml-splitter
```

### Verify Checksums
```console
# Download checksums
curl -L -o checksums.txt https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/checksums.txt

# Verify binary (Linux example)
sha256sum -c checksums.txt --ignore-missing
```

### Install to PATH
```console
sudo mv k8s-yaml-splitter /usr/local/bin/
```

## Usage

```console
Usage: k8s-yaml-splitter -f <input-file> <output-dir>
   or: k8s-yaml-splitter -f - <output-dir>  (read from stdin)
Options:
  -continue-on-error
        Continue processing on individual document errors (default true)
  -d    Create output directory if it doesn't exist
  -dry-run
        Dry run mode
  -exclude string
        Comma-separated list of resource kinds to exclude (e.g., 'Secret,ConfigMap')
  -f string
        Input file (use '-' for stdin)
  -include string
        Comma-separated list of resource kinds to include (e.g., 'Deployment,Service')
  -namespace-dirs
        Organize output files by namespace directories
  -o string
        Output format (yaml or json) (default "yaml")
  -s    Sort keys in output
  -version
        Show version information
```

## Examples:

#### Basic usage:
```console
k8s-yaml-splitter -f testdata/multi-namespace.yaml output-dir
```

#### Read from stdin (kubectl-style):
```console
kubectl get all -o yaml | k8s-yaml-splitter -f - output-dir
```

#### With sorted YAML output:
```console
k8s-yaml-splitter -f testdata/filter-test.yaml -s output-dir
```

#### JSON output (always sorted):
```console
k8s-yaml-splitter -f testdata/multi-namespace.yaml -o json output-dir
```

#### Create output directory automatically:
```console
k8s-yaml-splitter -f testdata/filter-test.yaml -d output-dir
```

#### Organize by namespace directories:
```console
k8s-yaml-splitter -f testdata/multi-namespace.yaml -namespace-dirs -d output-dir
```

#### Filter specific resource types:
```console
k8s-yaml-splitter -f testdata/filter-test.yaml -include "Deployment,Service" output-dir
```

#### Exclude sensitive resources:
```console
k8s-yaml-splitter -f testdata/filter-test.yaml -exclude "Secret,ConfigMap" output-dir
```

#### Dry run:
```console
k8s-yaml-splitter -f testdata/multi-namespace.yaml -dry-run output-dir
```

#### Combined options:
```console
k8s-yaml-splitter -f testdata/multi-namespace.yaml -s -o json -namespace-dirs -d output-dir
```

## Real-world Examples:

#### Split kubectl output:
```console
kubectl get all -n production -o yaml | k8s-yaml-splitter -f - -namespace-dirs ./manifests
```

#### Split Helm templates:
```console
helm template myapp ./chart | k8s-yaml-splitter -f - -exclude Secret ./output
```

#### Split Kustomize output:
```console
kustomize build ./overlay | k8s-yaml-splitter -f - -include Deployment,Service ./deploy
```

## Development

### Building from Source
```console
git clone https://github.com/ohauer/k8s-yaml-splitter.git
cd k8s-yaml-splitter
make build
```

### Testing
Run the test suite to verify functionality:
```console
make test
```

### Available Make Targets
```console
make help
```

## Example Output:

#### Dry Run:
```console
# k8s-yaml-splitter -f testdata/multi-namespace.yaml -dry-run output/
Found! type: Namespace | apiVersion: v1 | name: frontend | namespace:
==> DryRun: Writing output/Namespace-frontend.yaml
Found! type: Namespace | apiVersion: v1 | name: backend | namespace:
==> DryRun: Writing output/Namespace-backend.yaml
Found! type: Deployment | apiVersion: apps/v1 | name: web-server | namespace: frontend
==> DryRun: Writing output/Deployment-frontend-web-server.yaml
```

#### Normal Run:
```console
# k8s-yaml-splitter -f testdata/filter-test.yaml output/
Found! type: Deployment | apiVersion: apps/v1 | name: app-deployment | namespace: production
* Writing output/Deployment-production-app-deployment.yaml
* Wrote 158 bytes to output/Deployment-production-app-deployment.yaml
Found! type: Service | apiVersion: v1 | name: app-service | namespace: production
* Writing output/Service-production-app-service.yaml
* Wrote 142 bytes to output/Service-production-app-service.yaml
Found! type: ConfigMap | apiVersion: v1 | name: app-config | namespace: production
* Writing output/ConfigMap-production-app-config.yaml
* Wrote 110 bytes to output/ConfigMap-production-app-config.yaml

=== Processing Summary ===
Total documents: 5
Processed: 5
Skipped: 0
Errors: 0
```

#### With Namespace Directories:
```console
# k8s-yaml-splitter -f testdata/multi-namespace.yaml -namespace-dirs output/
Found! type: Namespace | apiVersion: v1 | name: frontend | namespace:
* Writing output/cluster-scoped/Namespace-frontend.yaml
* Wrote 145 bytes to output/cluster-scoped/Namespace-frontend.yaml
Found! type: Deployment | apiVersion: apps/v1 | name: web-server | namespace: frontend
* Writing output/frontend/Deployment-web-server.yaml
* Wrote 359 bytes to output/frontend/Deployment-web-server.yaml
Found! type: Service | apiVersion: v1 | name: web-service | namespace: frontend
* Writing output/frontend/Service-web-service.yaml
* Wrote 177 bytes to output/frontend/Service-web-service.yaml
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

This project is licensed under the same terms as the original k8s-yaml-splitter.

## Acknowledgments

Originally based on [k8s-yaml-splitter](https://github.com/latchmihay/k8s-yaml-splitter) by latchmihay, then forked to [k8s-yaml-splitter](https://github.com/mintel/k8s-yaml-splitter) by Mintel. Both upstream repositories have not responded to pull requests, so this repository has been disconnected from the fork and significantly enhanced and rewritten for production use.

