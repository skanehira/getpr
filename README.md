# getpr
Get GitHub's pull request URL.

## Usage

```sh
getpr - Get GitHub's Pull Request URL.

VERSION: 0.0.1

USAGE:
  $ getpr [OWNER/REPO] {commit id}

EXAMPLE:
  $ getpr getpr 737302e
  $ getpr getpr skanehira/getpr 737302e
```

## Installation
1. install getpr
   ```sh
   $ git clone https://github.com/skanehira/getpr
   $ cd getpr
   $ go install
   ```

2. please set GitHub token to `GITHUB_TOKEN` or `$HOME/.github_token`

## Author
skanehira
