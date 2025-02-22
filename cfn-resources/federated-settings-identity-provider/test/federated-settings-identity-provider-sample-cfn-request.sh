#!/usr/bin/env bash
# federated-setting-org-configs-sample-cfn-request.sh
#
# This tool generates text for a `cfn invoke` request json message.
#

set -o errexit
set -o nounset
set -o pipefail


function usage {
    echo "usage:$0 <FederationSettingsId> <ConnectedOrganizationId>"
    exit 1
}

if [ "$#" -ne 2 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

federationSettingsId="${1}"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$connectedOrganizationId" \
   '.desiredResourceState.FederationSettingsId?|=$FederationSettingsId | .desiredResourceState.ApiKeys.PublicKey?|=$pubkey | .desiredResourceState.ApiKeys.PrivateKey?|=$pvtkey ' \
     "$(dirname "$0")/federated-settings-identity-provider.sample-cfn-request.json"
