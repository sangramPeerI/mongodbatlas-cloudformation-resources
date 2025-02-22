#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#


function usage {
    echo "usage:$0 <project_name>"
}

if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

rm -rf inputs
mkdir inputs

#create apikey
org_id="$ATLAS_ORG_ID"

api_key_id=$(atlas organizations apikeys create --desc "cfn-boto-key-${CFN_TEST_TAG}" --role ORG_MEMBER --output json | jq -r '.id')

#create team
team_name="cfn-boto-team-${CFN_TEST_TAG}"
user_name=$(atlas organizations users list --output json | jq -r '.results[0].emailAddress')
echo "${user_name}"
team_id=$(atlas teams describe --name $team_name --output json | jq -r ".id")
if [ -z "$team_id" ]; then
  team_id=$(atlas team create "${team_name}" --username "${user_name}" --orgId "${org_id}" --output json | jq -r '.id')
fi


name="${1}"

projectId=$(atlas projects list --output json | jq --arg NAME "${name}" -r '.results[] | select(.name==$NAME) | .id')
if [ -z "$projectId" ]; then
      echo -e "No project found with \"${name}"
else
      echo -e "project found with ${name} and id ${projectId}, deleting"
      atlas projects delete "${projectId}" --force
fi

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "$name" \
   --arg key_id "$api_key_id" \
   --arg team_id "$team_id" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name | .ProjectApiKeys[0].Key?|=$key_id | .ProjectTeams[0].TeamId?|=$team_id' \
   "$(dirname "$0")/inputs_1_create.template.json" > "inputs/inputs_1_create.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "${name}- more B@d chars !@(!(@====*** ;;::" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name' \
   "$(dirname "$0")/inputs_1_invalid.template.json" > "inputs/inputs_1_invalid.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "${name}" \
   --arg key_id "$api_key_id" \
   --arg team_id "$team_id" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name | .ProjectApiKeys[0].Key?|=$key_id | .ProjectTeams[0].TeamId?|=$team_id'\
   "$(dirname "$0")/inputs_1_update.template.json" > "inputs/inputs_1_update.json"

ls -l inputs

echo "TODO: Delete the team and api_key created above"
