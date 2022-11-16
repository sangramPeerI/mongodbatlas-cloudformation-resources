#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#

set -o errexit
set -o nounset
set -o pipefail

function usage {
    echo "usage:$0 <project_name>"
}

if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

rm -rf inputs
mkdir inputs

projectName="$1"
ATLAS_PROJECT_ID=$(mongocli iam projects list --output json | jq --arg NAME "${projectName}" -r '.results[] | select(.name==$NAME) | .id')
if [ -z "$ATLAS_PROJECT_ID" ]; then
    ATLAS_PROJECT_ID=$(mongocli iam projects create "${projectName}" --output=json | jq -r '.id')

    echo -e "Created project \"${projectName}\" with id: ${ATLAS_PROJECT_ID}\n"
else
    echo -e "FOUND project \"${projectName}\" with id: ${ATLAS_PROJECT_ID}\n"
fi

ATLAS_REGION_NAME="US_EAST_1"


jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg projectId "$ATLAS_PROJECT_ID" \
   --arg regionName "$ATLAS_REGION_NAME" \
   '.ProjectId?|=$projectId | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .RegionName?|=$regionName' \
   "$(dirname "$0")/inputs_1_create.template.json" > "inputs/inputs_1_create.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg projectId "$ATLAS_PROJECT_ID" \
   --arg regionName "$ATLAS_REGION_NAME" \
   '.ProjectId?|=$projectId | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .RegionName?|=$regionName' \
   "$(dirname "$0")/inputs_1_invalid.template.json" > "inputs/inputs_1_invalid.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg projectId "$ATLAS_PROJECT_ID" \
   --arg regionName "$ATLAS_REGION_NAME" \
   '.ProjectId?|=$projectId | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .RegionName?|=$regionName' \
   "$(dirname "$0")/inputs_1_update.template.json" > "inputs/inputs_1_update.json"

ls -l inputs

echo "TODO: Delete the team and api_key created above"
