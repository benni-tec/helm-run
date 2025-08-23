Set-Location $Env:HELM_PLUGIN_DIR
git fetch --tags
git checkout --quiet main
git checkout --quiet origin/release/$(git describe --tags $(git rev-list --tags --max-count=1))
