#!/bin/bash

LATEST_TAG=$(git tag | tail - )

if [[ $GITHUB_TOKEN == "" ]]; then
echo "Github token not set. Aborting."
fi

if [[ $1 != $LATEST_TAG ]]; then
echo "Latest tag does not match build tag"
exit 1
fi

git checkout master
git push --tags

go build

if ![ -f gaea-server ]; then
echo "Build failed"
exit 1
fi

tar -cvzf gaea-server-linux-amd64.tar.gz gaea-server

github-release release \
    --user btburke \
    --repo gaea-server \
    --tag $LATEST_TAG \
    --name "$LATEST_TAG" \
    --description "API server for the GAEA website"

github-release upload \
    --user btburke \
    --repo gaea-server \
    --tag $LATEST_TAG \
    --name "gaea-server-linux-amd64.tar.gz" \
    --file gaea-server-linux-amd64.tar.gz

rm gaea-server-linux-amd64.tar.gz

echo "Release $LATEST_TAG complete and uploaded to Github"