#!/bin/bash

LATEST_TAG=$(git tag | tail -n 1)

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

if [ -f gaea-server ]; then

tar -cvzf gaea-server-linux-amd64.tar.gz gaea-server
tar -cvzf sql.tar.gz docker/db/schema.sql /docker/db/sql/

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

github-release upload \
    --user btburke \
    --repo gaea-server \
    --tag $LATEST_TAG \
    --name "sql.tar.gz" \
    --file sql.tar.gz


rm gaea-server-linux-amd64.tar.gz
rm sql.tar.gz

else
echo "Build failed"
exit 1
fi


echo "Release $LATEST_TAG complete and uploaded to Github"