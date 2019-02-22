# Gitup

Ensure a list of repository is cloned in the right place.

This script is useful when recreating your workspace from scratch so you don't have to clone all repositories yourself,
one by one.

## Usage

### YAML configuration file

Create a configuration file in `~/.gituprc` with the following format:

```yaml
repositories:
  - url: git@github.com/your_name/a_repository
    path: ~/Workspace/a_repository
  - url: git@github.com/an_organization/another_repository
    path: ~/Workspace/another_repository
```

> :warning: Be aware the script doesn't handle SSH authentication so you have to set your SSH Agent before calling
> the script:
>
> ```
> ssh-agent ~/.ssh/id_rsa
> ```
