# racelogctl
Command line interface to interact with racelogger service


## Create release notes

```console
gh api --method POST -H "Accept: application/vnd.github+json" /repos/mpapenbr/racelogctl/releases/generate-notes -f tag_name='v0.6.0'  | jq -r '.body'
```