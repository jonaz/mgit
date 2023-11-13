# mgit
Multi repo edit/commit/pr tool

## Install

```
go install github.com/jonaz/mgit@latest
```

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



### example to clone repos containing content
```
mgit -c https://git.domain clone --has-file config.json --content-regexp '"team": "mycoolteam"'
mgit -c https://git.domain playbook generate
# edit the playbook
mgit -c https://git.domain playbook run playbook.yml
mgit -c https://git.domain playbook pr --mode api playbook.yml
```

## docs

### actions

| field | description|
| --- | ----------- |
| command | runs command in shell. cannot be used with regexp. If used together with contentRegexp|fileRegexp|pathRegexp it invoke the command per file found with {{.FilePath}} | 
| regexp | |
| with | what to replace matched regexp with. capture groups example $1 is supported |
| fileRegexp | only change in files where the filename matches the regexp |
| pathRegexp | only change files in matching path. Includes the full path from the repo root and not only the filename |
| contentRegexp | only change in files which match this regexp. |

