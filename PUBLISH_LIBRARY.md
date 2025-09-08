## Publish Library
1. git tag "vx.x.x"
2. git push --tag

## go get private repos using SSH key auth instead of password auth
1. add export GOPRIVATE=github.com/{{your usename}}
2. add export GITHUB_TOKEN=your-access-token-generated-from-github
3. git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
4. cat .gitconfig


Look something like this<br/>
[url "https://your-access-token-generated-from-github:x-oauth-basic@github.com/"]

        insteadOf = https://github.com/
