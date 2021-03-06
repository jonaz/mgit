# mgit
Multi repo edit/commit/pr tool

## Example usage

there are help for all command and subcommands

```
mgit --help
```
example clone and replace

```
mgit --dir=tmp --bitbucket-url=https://bitbucket-server clone --whitelist my-cool-repo,my-cool-repo2
mgit --dir=tmp replace --with asdfasdf\$1 --regexp "gograce.New(.*)" --file-regexp=go$
```

example command

```
mgit --dir=tmp command go mod tidy
```

### Run using playbook.yml 

Example using a playbook to do multiple things to a repo
playbook.yml:

```
tasks:
- repos:
    - ssh://git@git.asdf.com/repo1.git
    - ssh://git@git.asdf.com/repo2.git
    - ssh://git@git.asdf.com/repo3.git
  targetBranch: update-gograce
  commitMessage: |
    Update gograce to support sleep before shutting down

    set sleep to 5 sec which we have tested and works fine with k8s and traefik to get rid of 502's when we deploy
  actions:
    - command: go get -u github.com/jonaz/gograce
    - command: go mod tidy
    - regexp: gograce.NewServerWithTimeout\((.*)\)
      with: "gograce.NewServerWithTimeout($1, 5 * time.Second)"
      fileRegexp: "go$"
    - command: go fmt ./...
```

run with `mgit playbook run playbook.yml`
open multiple PRs when done: `mgit --bitbucket-url=https://bitbucketserver.com pr`
