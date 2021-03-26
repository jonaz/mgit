# mgit
Multi repo edit/commit/pr tool

### Example usage

```
mgit --dir=tmp --bitbucket-url=https://bitbucket-server clone --whitelist my-cool-repo,my-cool-repo2
mgit --dir=tmp --bitbucket-url=https://bitbucket-server replace --with asdfasdf\$1 --regexp "gograce.New(.*)" --file-regexp=go$
```
