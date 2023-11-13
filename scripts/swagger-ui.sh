#!/usr/bin/env bash
REPO="https://github.com/swagger-api/swagger-ui"
DIR=$(basename "$REPO")

rm -rf "$DIR"
git clone "$REPO" --depth=1
mkdir swaggerui -p
cp "$DIR/dist" swaggerui -r
grep -Rl "petstore.swagger.io" swaggerui | xargs sed -i 's/petstore.swagger.io\/v2\/swagger.json/shopping-list.apachejuice.dev\/openapi\/swagger.yaml/g'
cp spec/swagger.yaml "swaggerui/dist"
rm "$DIR" -rf
