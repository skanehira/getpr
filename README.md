# getpr
Get GitHub Enterprise's pull request URL.

![](https://i.imgur.com/VrXQw15.gif)

## Usage

```sh
getpr - Get GitHub Enterprise's Pull Request URL.

VERSION: 0.0.1

USAGE:
  $ getpr [OWNER/REPO] {commit id}

EXAMPLE:
  $ getpr 02b3cb3
  $ getpr skanehira/getpr 02b3cb3
```

## Installation
1. install getpr
   ```sh
   $ git clone https://github.com/skanehira/getpr
   $ cd getpr
   $ go install
   ```

2. please set GitHub token to `GITHUB_TOKEN` or `$HOME/.github_token`
3. please set GitHub Enterprise Graphql API Endpoint to `GITHUB_ENDPOINT` (ex: `https://git.hoge.com/api/graphql`)

## Author
skanehira
