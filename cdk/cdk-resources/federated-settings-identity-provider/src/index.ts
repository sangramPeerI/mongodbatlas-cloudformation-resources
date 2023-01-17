// Generated by cdk-import
import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';

/**
 * Returns, adds, edits, and removes federation-related features such as role mappings and connected organization configurations.
 *
 * @schema CfnFederatedSettingsIdentityProviderProps
 */
export interface CfnFederatedSettingsIdentityProviderProps {
  /**
   * @schema CfnFederatedSettingsIdentityProviderProps#ApiKeys
   */
  readonly apiKeys?: ApiKeyDefinition;

  /**
   * Flag that indicates whether the identity provider has SSO debug enabled.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#SsoDebugEnabled
   */
  readonly ssoDebugEnabled?: boolean;

  /**
   * List that contains the domains associated with the identity provider.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#AssociatedDomains
   */
  readonly associatedDomains?: string[];

  /**
   * String enum that indicates whether the identity provider is active.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#Status
   */
  readonly status?: string;

  /**
   * Unique string that identifies the issuer of the SAML Assertion.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#IssuerUri
   */
  readonly issuerUri?: string;

  /**
   * SAML Authentication Request Protocol HTTP method binding (POST or REDIRECT) that Federated Authentication uses to send the authentication request.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#RequestBinding
   */
  readonly requestBinding?: string;

  /**
   * Signature algorithm that Federated Authentication uses to encrypt the identity provider signature.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#ResponseSignatureAlgorithm
   */
  readonly responseSignatureAlgorithm?: string;

  /**
   * URL that points to the receiver of the SAML authentication request.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#SsoUrl
   */
  readonly ssoUrl?: string;

  /**
   * Human-readable label that identifies the IdP.
   *
   * @schema CfnFederatedSettingsIdentityProviderProps#DisplayName
   */
  readonly displayName?: string;

}

/**
 * Converts an object of type 'CfnFederatedSettingsIdentityProviderProps' to JSON representation.
 */
/* eslint-disable max-len, quote-props */
export function toJson_CfnFederatedSettingsIdentityProviderProps(obj: CfnFederatedSettingsIdentityProviderProps | undefined): Record<string, any> | undefined {
  if (obj === undefined) { return undefined; }
  const result = {
    'ApiKeys': toJson_ApiKeyDefinition(obj.apiKeys),
    'SsoDebugEnabled': obj.ssoDebugEnabled,
    'AssociatedDomains': obj.associatedDomains?.map(y => y),
    'Status': obj.status,
    'IssuerUri': obj.issuerUri,
    'RequestBinding': obj.requestBinding,
    'ResponseSignatureAlgorithm': obj.responseSignatureAlgorithm,
    'SsoUrl': obj.ssoUrl,
    'DisplayName': obj.displayName,
  };
  // filter undefined values
  return Object.entries(result).reduce((r, i) => (i[1] === undefined) ? r : ({ ...r, [i[0]]: i[1] }), {});
}
/* eslint-enable max-len, quote-props */

/**
 * @schema ApiKeyDefinition
 */
export interface ApiKeyDefinition {
  /**
   * @schema ApiKeyDefinition#PrivateKey
   */
  readonly privateKey?: string;

  /**
   * @schema ApiKeyDefinition#PublicKey
   */
  readonly publicKey?: string;

}

/**
 * Converts an object of type 'ApiKeyDefinition' to JSON representation.
 */
/* eslint-disable max-len, quote-props */
export function toJson_ApiKeyDefinition(obj: ApiKeyDefinition | undefined): Record<string, any> | undefined {
  if (obj === undefined) { return undefined; }
  const result = {
    'PrivateKey': obj.privateKey,
    'PublicKey': obj.publicKey,
  };
  // filter undefined values
  return Object.entries(result).reduce((r, i) => (i[1] === undefined) ? r : ({ ...r, [i[0]]: i[1] }), {});
}
/* eslint-enable max-len, quote-props */


/**
 * A CloudFormation `MongoDB::Atlas::FederatedSettingsIdentityProvider`
 *
 * @cloudformationResource MongoDB::Atlas::FederatedSettingsIdentityProvider
 * @stability external
 */
export class CfnFederatedSettingsIdentityProvider extends cdk.CfnResource {
  /**
  * The CloudFormation resource type name for this resource class.
  */
  public static readonly CFN_RESOURCE_TYPE_NAME = 'MongoDB::Atlas::FederatedSettingsIdentityProvider';

  /**
   * Resource props.
   */
  public readonly props: CfnFederatedSettingsIdentityProviderProps;

  /**
   * Attribute `MongoDB::Atlas::FederatedSettingsIdentityProvider.FederationSettingsId`
   */
  public readonly attrFederationSettingsId: string;
  /**
   * Attribute `MongoDB::Atlas::FederatedSettingsIdentityProvider.OktaIdpId`
   */
  public readonly attrOktaIdpId: string;

  /**
   * Create a new `MongoDB::Atlas::FederatedSettingsIdentityProvider`.
   *
   * @param scope - scope in which this resource is defined
   * @param id    - scoped id of the resource
   * @param props - resource properties
   */
  constructor(scope: constructs.Construct, id: string, props: CfnFederatedSettingsIdentityProviderProps) {
    super(scope, id, { type: CfnFederatedSettingsIdentityProvider.CFN_RESOURCE_TYPE_NAME, properties: toJson_CfnFederatedSettingsIdentityProviderProps(props)! });

    this.props = props;

    this.attrFederationSettingsId = cdk.Token.asString(this.getAtt('FederationSettingsId'));
    this.attrOktaIdpId = cdk.Token.asString(this.getAtt('OktaIdpId'));
  }
}