# k8s values updater

Changes certain values on your deployment values, updates the tag and pullRequestID key and creates a git commit

- The add to know hosts and push are not working at the moment
- by default it will search every file with the name `values.yaml` inside evey folder on `deploy/*`
- Every option is optional, if you run with no options it will only indent your values files and remove empty keys

## Options

The main options of the bump command are:
- `--tag` Set the `chart.image.tag` value
- `--pr-id` Set the `chart.image.pullRequestId` value
- `--file-names` Can be a comma separetaed value like `qa-values.yaml,live-values.yaml`
- `--dir-path` Set the root file that the cli will search, accepts wildcard
- `--is-root` If set will assume there is not a root chart key, you are not using any root chart, will change the values inside `image.tag` for example
- `--replace-with` You can pass a yaml and the cli will try to merge it

Here`s the help
```
» k8s-values-updater
Bump value on file

Usage:
  k8s-values-updater bump [flags]

Flags:
  -b, --branch string                Branch (Will use if set $CIRCLE_BRANCH)
      --chart-name string            The name of the subchart
      --dir-path string              File Path (default "deploy/*")
  -e, --email string                 Email that will commit (default "k8s-values-updater@mailinator.com")
      --file-names string            File Path (default "values.yaml")
      --github-access-token string   Github Acccess Token (Will use if set $GITHUB_ACCESS_TOKEN)
  -h, --help                         help for bump
      --is-root                      If set will define that the values to be changed has no subchart
  -o, --org string                   Organization (Will use if set $CIRCLE_PROJECT_USERNAME)
  -p, --pr-id string                 Pull Request ID
      --remote-name string            (default "origin")
      --replace-with string          If passed will try to merge this value with the values yaml
  -r, --repo string                  Repository (Will use if set $CIRCLE_PROJECT_REPONAME)
  -t, --tag string                   Image Tag
  -w, --workdir string               Workdir (default ".")

Global Flags:
      --config string      config file (default is $HOME/.config.yaml)
      --dry-run            Dry Run
      --log-level string   Log level (debug, info, warn, error, fatal, panic (default "warning")
  -v, --verbose            Dry Run
```

## Local Usage

## Development

## Circle

```
# ...
workflows:
  build-and-test:
    jobs:
      - test:
          context: <CONTEXT_WITH_GITHUB_TOKEN>
# ...
jobs:
  # Test
  test:
    docker:
      - image: quay.io/dafiti/k8s-values-updater:master
    steps:
      - setup_remote_docker
      - run:
          name: Try to commit and push
          command: |
            /k8s-values-updater bump -e <some-email> -t $(echo ${CIRCLE_SHA1} | head -c7)

```

...
