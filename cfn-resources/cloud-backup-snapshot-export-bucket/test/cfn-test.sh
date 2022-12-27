#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#

function usage {
    echo "Creates a new private endpoint role for the test"
}

roleName="mongodb-test-cloud-backup-export-bucket-role-US_EAST_2"
policyName="atlas-kms-role-policy-US_EAST_2"

awsRoleID=$(aws iam get-role --role-name "${roleName}" | jq --arg roleName "${roleName}" -r '.Role | select(.RoleName==$roleName) |.RoleId')
if [ -z "$awsRoleID" ]; then
    echo -e "No role found, hence creating the role\n $awsRoleID"
else
    echo -e "Role found $awsRoleID"
fi