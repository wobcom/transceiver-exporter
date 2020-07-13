#!/bin/bash

set -e

echo "Publishing artifacts"

if [ -z $DEPLOY_TOKEN ]; then
    echo "DEPLOY_TOKEN not set!"
    exit 42
fi

token="$DEPLOY_TOKEN"
gitlab="https://gitlab.com"
api="$gitlab/api/v4"


echo "Uploading the binary to $gitlab"
out=$(curl -f \
	   --request POST \
           --header "PRIVATE-TOKEN: $token" \
           --form "file=@$CI_PROJECT_DIR/transceiver-exporter" \
	   "$api/projects/$CI_PROJECT_ID/uploads")


echo "Response from gitlab is:"
echo "$out"
url=$(echo "$out" | jq -r '.full_path')

body=$(cat <<JSON
{
  "tag_name": "$CI_COMMIT_TAG",
  "name": "$ref",
  "assets": {
    "links": [
      { "name": "transceiver-exporter",
        "url": "$gitlab$url",
        "filepath": "/binaries/transceiver-exporter"
      }
    ]
  }
}
JSON
)

echo "Using the following body..."
echo "$body" 

echo "... creating a release"
curl -f \
     -o - \
     --header 'Content-Type: application/json' \
     --header "PRIVATE-TOKEN: $token" \
     --data "$body" \
     --request POST \
     "$api/projects/$CI_PROJECT_ID/releases"

