# Github Repository Lister

Reads a Github PAT from `token.txt`, writes a list of repositories to `.env`:
```
GITHUB_REPOS=${org-name}/${repo-name},
```

## Usage:

```
$ echo ${GH_PAT} > token.txt
$ ./repo-lister -org=org-name
```