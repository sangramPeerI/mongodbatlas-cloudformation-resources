#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#

#set -x
function usage {
    echo "usage:$0 <project_name>"
    echo "Creates a new encryption key for the the project "
}



if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi
rm -rf inputs
mkdir inputs

projectName="${1}"
projectId=$(atlas projects list --output json | jq --arg NAME "${projectName}" -r '.results[] | select(.name==$NAME) | .id')
if [ -z "$projectId" ]; then
    projectId=$(atlas projects create "${projectName}" --output=json | jq -r '.id')

    echo -e "Created project \"${projectName}\" with id: ${projectId}\n"
else
    echo -e "FOUND project \"${projectName}\" with id: ${projectId}\n"
fi

echo "Check if a project is created $projectId"
export MCLI_PROJECT_ID=$projectId

keyRegion=$AWS_DEFAULT_REGION
if [ -z "$keyRegion" ]; then
keyRegion=$(aws configure get region)
fi
keyRegion=$(echo "$keyRegion" | sed -e "s/-/_/g")
keyRegion=$(echo "$keyRegion" | tr '[:lower:]' '[:upper:]')
echo "$keyRegion"

roleName="mongodb-test-enc-role-${keyRegion}"
policyName="atlas-kms-role-policy-${keyRegion}"

echo "roleName: ${roleName} , policyName: ${policyName}"

echo "--------------------------------create key and key policy document starts ----------------------------"\n

keyARN=$(aws kms create-key | jq '.KeyMetadata|.Arn')
prefix='{ "Version": "2012-10-17", "Statement": ['
echo "--------------------------------printing key  starts ----------------------------"\n
echo "$keyARN"
cleanedkeyARN=$(echo ${keyARN} | sed 's/"//g')
echo $cleanedkeyARN
echo "--------------------------------printing key  ends ----------------------------"\n

policyContent=$(jq --arg cleanedkeyARN $cleanedkeyARN '.Statement[0]|.Resource[0]?|=$cleanedkeyARN' "$(dirname "$0")/key-policy-template.json" )
suffix=']}'
policyDocument="${prefix} ${policyContent} ${suffix}"
echo $policyDocument > $(dirname "$0")/policy.json

policyContent=$(jq '.Statement[0].Resource[0]' "$(dirname "$0")/policy.json" )
echo "$policyContent"
keyID=$(echo ${policyContent##*/})
cleanedKeyID=$(echo "${keyID}" | sed 's/"//g')
echo $cleanedKeyID

echo "--------------------------------create key and key policy document policy document ends ----------------------------"\n


echo $policyDocument
echo "--------------------------------policy document finished ----------------------------"\n


roleID=$(atlas cloudProviders accessRoles aws create --output json | jq -r '.roleId')
echo "--------------------------------Mongo CLI Role creation ends ----------------------------"\n


atlasAWSAccountArn=$(atlas cloudProviders accessRoles  list --output json | jq --arg roleID "${roleID}" -r '.awsIamRoles[] |select(.roleId |test( $roleID)) |.atlasAWSAccountArn')
atlasAssumedRoleExternalId=$(atlas cloudProviders accessRoles  list --output json | jq --arg roleID "${roleID}" -r '.awsIamRoles[] |select(.roleId |test( $roleID)) |.atlasAssumedRoleExternalId')
jq --arg atlasAssumedRoleExternalId "$atlasAssumedRoleExternalId" \
   --arg atlasAWSAccountArn "$atlasAWSAccountArn" \
  '.Statement[0].Principal.AWS?|=$atlasAWSAccountArn | .Statement[0].Condition.StringEquals["sts:ExternalId"]?|=$atlasAssumedRoleExternalId' "$(dirname "$0")/role-policy-template.json" >"$(dirname "$0")/add-policy.json"
echo cat add-policy.json
echo "--------------------------------AWS Role creation ends ----------------------------"\n


awsRoleID=$(aws iam get-role --role-name "${roleName}" | jq --arg roleName "${roleName}" -r '.Role | select(.RoleName==$roleName) |.RoleId')
if [ -z "$awsRoleID" ]; then
    awsRoleID=$(aws iam create-role --role-name "${roleName}" --assume-role-policy-document file://$(dirname "$0")/add-policy.json | jq --arg roleName "${roleName}" -r '.Role | select(.RoleName==$roleName) |.RoleId')
    echo -e "No role found, hence creating the role. Created id: ${awsRoleID}\n"
else
    aws iam delete-role-policy --role-name "${roleName}" --policy-name "${policyName}"
    aws iam delete-role --role-name "${roleName}"
 awsRoleID=$(aws iam create-role --role-name "${roleName}" --assume-role-policy-document file://$(dirname "$0")/add-policy.json | jq --arg roleName "${roleName}" -r '.Role | select(.RoleName==$roleName) |.RoleId')
    echo -e "FOUND id: ${awsRoleID}\n"
fi
echo "--------------------------------AWS Role creation ends ----------------------------"\n


awsArn=$(aws iam get-role --role-name "${roleName}" | jq --arg roleName "${roleName}" -r '.Role | select(.RoleName==$roleName) |.Arn')

aws iam put-role-policy   --role-name "${roleName}"   --policy-name "${policyName}"   --policy-document file://$(dirname "$0")/policy.json
echo "--------------------------------attach mongodb  Role to AWS Role ends ----------------------------"\n

awsArne=$(echo "${awsArn}" | sed 's/"//g')
# shellcheck disable=SC2086
#TODO Needs change to while loop using get operation
sleep 65

atlas cloudProviders accessRoles aws authorize ${roleID} --iamAssumedRoleArn ${awsArne}
echo "--------------------------------authorize mongodb  Role ends ----------------------------"\n

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg projectId "$projectId" \
   --arg KMS_KEY "$cleanedKeyID" \
   --arg KMS_ROLE "${roleID}" \
   --arg region "$keyRegion" \
   '.AwsKms.CustomerMasterKeyID?|=$KMS_KEY | .AwsKms.RoleID?|=$KMS_ROLE | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .ProjectId?|=$projectId | .AwsKms.Region?|=$region ' \
   "$(dirname "$0")/inputs_1_create.template.json" > "inputs/inputs_1_create.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg KMS_KEY "$cleanedKeyID" \
   --arg KMS_ROLE "${roleID}" \
   --arg projectId "$projectId" \
   --arg region "$keyRegion" \
    '.AwsKms.CustomerMasterKeyID?|=$KMS_KEY | .AwsKms.RoleID?|=$KMS_ROLE | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .ProjectId?|=$projectId | .AwsKms.Region?|=$region' \
   "$(dirname "$0")/inputs_1_invalid.template.json" > "inputs/inputs_1_invalid.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg projectId "$projectId" \
   --arg KMS_KEY "$cleanedKeyID" \
   --arg KMS_ROLE "${roleID}" \
   --arg region "$keyRegion" \
   '.AwsKms.CustomerMasterKeyID?|=$KMS_KEY | .AwsKms.RoleID?|=$KMS_ROLE | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .ProjectId?|=$projectId | .AwsKms.Region?|=$region ' \
   "$(dirname "$0")/inputs_1_update.template.json" > "inputs/inputs_1_update.json"
ls -l inputs
#mongocli iam projects delete "${projectId}" --force



#mongocli atlas cloudProviders accessRoles aws authorize 63721b924ad9a46eeef105ae --iamAssumedRoleArn "arn:aws:iam::816546967292:role/mongodb-test-enc-role"