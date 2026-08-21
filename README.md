# gd-demo-4-gates — step 4: guard rails

Fourth of the five graduated demo repositories for
[git-deploy-operator](https://github.com/fsamin/git-deploy-operator).

**What this repository proves:** a commit that fails a blocking gate is never
released. The application keeps serving the previous commit, and the report
stays readable either way.

The page prints an order total computed by `price.go`. The sample basket is
worth exactly the discount threshold, which is the boundary a regression is most
likely to move — so a broken business rule is visible as a wrong total on the
page. The demo is that you never get to see it.

## Demo

```sh
git-deploy -n demo init --name gates
git-deploy edit                     # add the four gates (or use the UI)
```

```yaml
spec:
  gates:
    - type: test
    - type: quality
    - type: vulncheck
    - type: imageScan
```

Four green tabs in the log panel. Then push the regression onto the tracked
branch:

```sh
git cherry-pick failing-test        # one character: `<` becomes `<=`
git push
```

Within one poll interval the phase goes `Failed`, the *test* tab opens on the
failure by itself — and the page still serves the previous total. The proof is
the absence of a release, not the phase:

```sh
kubectl get appreleases -n demo     # no release for the new commit
```

Then revert, and the same pipeline releases the good commit again.

## One branch per gate

| Branch | What it breaks | test | quality | vulncheck | imageScan |
|---|---|:--:|:--:|:--:|:--:|
| `main` | nothing | ✓ | ✓ | ✓ | ✓ |
| `failing-test` | the discount boundary: `<` becomes `<=`, so an order landing exactly on the threshold loses its discount | ✗ | ✓ | ✓ | ✓ |
| `lint-error` | leaves a dead function behind after extracting a helper | ✓ | ✗ | ✓ | ✓ |
| `vuln-dep` | parses the sample basket with `gopkg.in/yaml.v2@v2.2.2` | ✓ | ✓ | ✗ | ✗ |
| `image-cve` | ships the binary on `debian:11-slim` instead of `scratch` | ✓ | ✓ | ✓ | ✗ |

`vuln-dep` deliberately trips **two** gates, and the two reports are worth
reading side by side: `govulncheck` finds three vulnerabilities and names the
call that reaches them (`gd.loadBasket calls yaml.Unmarshal`), while `trivy`
finds one and names the module version sitting in the image. Same dependency,
two angles — a reachability analysis and an inventory — which is why both gates
exist.

Every other branch fails exactly one gate, so a demo can pick the one it wants
to show.

## Why the builder image is `golang:1.26-alpine`

Because `imageScan` scans the **binary** as well as the OS packages, and a Go
binary carries its toolchain's standard library. Built with `golang:1.24`, this
application failed the image gate with 19 HIGH CVEs in `crypto/x509`, `net/http`
and friends — from a `FROM scratch` image with no OS package at all. The
vulnerabilities were real; they were just not this repository's code.

The floating minor tag is deliberate: it tracks the latest patch release, so the
gate stays green as fixes ship instead of going red the week a CVE is published.

## Locally

```sh
go test ./...          # the discount rule and its boundary
golangci-lint run      # clean
go run .               # http://localhost:8080
```
