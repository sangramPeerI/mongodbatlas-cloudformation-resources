#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#
set -x
echo "--------------------------------delete key and key policy document policy document starts ----------------------------"\n


projectName="${1}"
projectId=$(atlas projects list --output json | jq --arg NAME "${projectName}" -r '.results[] | select(.name==$NAME) | .id')
echo "Check if a project is created $projectId"
export MCLI_PROJECT_ID=$projectId

keyRegion=$AWS_DEFAULT_REGION
awsRegion=$AWS_DEFAULT_REGION
if [ -z "$keyRegion" ]; then
keyRegion=$(aws configure get region)
fi
keyRegion=$(echo "$keyRegion" | sed -e "s/-/_/g")
keyRegion=$(echo "$keyRegion" | tr '[:lower:]' '[:upper:]')
echo "$keyRegion"

roleName="mongodb-test-cloud-backup-export-bucket-role-${keyRegion}"
policyName="atlas-cloud-backup-export-bucket-S3-role-policy-${keyRegion}"

pwd
trustPolicy=$(jq '.Statement[0].Condition.StringEquals["sts:ExternalId"]' "add-policy.json" )
echo $trustPolicy
roleExternalID=$(echo ${trustPolicy##*/})
atlasAssumedRoleExternalID=$(echo "${roleExternalID}" | sed 's/"//g')
echo $atlasAssumedRoleExternalID

roleId=$(atlas cloudProviders accessRoles list --output json --projectId ${projectId} | jq --arg roleID "${atlasAssumedRoleExternalID}" -r '.awsIamRoles[] |select(.atlasAssumedRoleExternalId |test( $roleID)) |.roleId')
echo $roleId

atlas cloudProviders accessRoles aws deauthorize ${roleId} --projectId ${projectId} --force
echo "--------------------------------delete role starts ----------------------------"\n

aws iam delete-role-policy --role-name "$roleName" --policy-name "$policyName"
aws iam delete-role --role-name "$roleName"
echo "--------------------------------delete role ends ----------------------------"\n

bucketName="cloud-backup-snapshot-test-"${awsRegion}

aws s3 rb s3://${bucketName} --force
echo "--------------------------------delete bucket ends ----------------------------"\n
#mongocli iam projects delete "${projectId}" --force