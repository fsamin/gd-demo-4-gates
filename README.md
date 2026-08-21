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

| Branch | What it breaks | The gate that catches it |
|---|---|---|
| `main` | nothing — all four green | — |
| `failing-test` | the discount boundary: `<` becomes `<=`, so an order landing exactly on the threshold loses its discount | `test` |
| `lint-error` | leaves a dead function behind after a refactor | `quality` |
| `vuln-dep` | pulls a YAML parser with a known, *reachable* vulnerability | `vulncheck` |
| `image-cve` | ships the binary on a long-unpatched base image instead of `scratch` | `imageScan` |

Each one fails exactly one gate and passes the other three, so a demo can pick
the gate it wants to show.

## Locally

```sh
go test ./...          # the discount rule and its boundary
golangci-lint run      # clean
go run .               # http://localhost:8080
```
