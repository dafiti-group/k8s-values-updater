# k8s values updater

Changes certain values on your deployment values, updates the tag and pullRequestID key and creates a git commit

- The add to know hosts and push are not working at the moment
- by default it will search every file with the name `values.yaml` inside evey folder on `deploy/*`
- Every option is optional, if you run with no options it will only indent your values files and remove empty keys

## Options

The main options of the bum command are:
- `--tag` Set the `chart.image.tag` value
- `--pr-id` Set the `chart.image.pullRequestId` value
- `--file-names` Can be a comma separetaed value like `qa-values.yaml,live-values.yaml`
- `--dir-path` Set the root file that the cli will search, accepts wildcard
- `--is-root` If set will assume there is not a root chart key, you are not using any root chart, will change the values inside `image.tag` for example
- `--replace-with` You can pass a yaml and the cli will try to merge it

Here`s the help
```
Bump value on file

Usage:
  k8s-values-updater bump [flags]

Flags:
      --branch string         branch
      --chart-name string     The name of the subchart
      --dir-path string       File Path (default "deploy/*")
      --email string          Email that will commit (default "k8s-values-updater@mailinator.com")
      --file-names string     File Path (default "values.yaml")
  -h, --help                  help for bump
      --is-root               If set will define that the values to be changed has no subchart
      --pr-id string          Pull Request ID
      --refspec string        refspec
      --remote-name string     (default "origin")
      --replace-with string   If passed will try to merge this value with the values yaml
      --tag string            Image Tag

Global Flags:
      --config string   config file (default is $HOME/.config.yaml)
      --dry-run         Dry Run
  -v, --verbose         verbose output
```

## Local Usage

## Development

## Circle

```
  # Update PR ID
  update-pr-id:
    docker:
      # https://github.com/dafiti-group/k8s-values-updater
      - image: quay.io/dafiti/k8s-values-updater:latest
    steps:
      - setup_remote_docker
      - add_ssh_keys:
          fingerprints:
            - "<your-ssh-key-fingerprint>"
      - checkout
      - run:
          name: Update PR ID
          command: |
            # Will remove this in the future
            ssh-keyscan -H github.com > ~/.ssh/known_hosts
            #
            # Will update the key pullRequestID on every values file
            /k8s-values-updater bump \
              --pr-id $(echo ${CI_PULL_REQUEST##*/}) \
              --email <some-email> \
            # Will remove this in the future
            git push origin $CIRCLE_BRANCH
```

...
