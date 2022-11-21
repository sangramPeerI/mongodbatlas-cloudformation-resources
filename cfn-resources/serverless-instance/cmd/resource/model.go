// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	ApiKeys                      *ApiKeyDefinition                    `json:",omitempty"`
	ConnectionStrings            *ServerlessInstanceConnectionStrings `json:",omitempty"`
	ContinuousBackupEnabled      *bool                                `json:",omitempty"`
	CreateDate                   *string                              `json:",omitempty"`
	Id                           *string                              `json:",omitempty"`
	IncludeCount                 *bool                                `json:",omitempty"`
	ItemsPerPage                 *int                                 `json:",omitempty"`
	Links                        []Link                               `json:",omitempty"`
	MongoDBVersion               *string                              `json:",omitempty"`
	Name                         *string                              `json:",omitempty"`
	PageNum                      *int                                 `json:",omitempty"`
	ProjectID                    *string                              `json:",omitempty"`
	ProviderSettings             *ServerlessInstanceProviderSettings  `json:",omitempty"`
	StateName                    *string                              `json:",omitempty"`
	TerminationProtectionEnabled *bool                                `json:",omitempty"`
	TotalCount                   *float64                             `json:",omitempty"`
}

// ApiKeyDefinition is autogenerated from the json schema
type ApiKeyDefinition struct {
	PrivateKey *string `json:",omitempty"`
	PublicKey  *string `json:",omitempty"`
}

// ServerlessInstanceConnectionStrings is autogenerated from the json schema
type ServerlessInstanceConnectionStrings struct {
	PrivateEndpoint []ServerlessInstancePrivateEndpoint `json:",omitempty"`
	StandardSrv     *string                             `json:",omitempty"`
}

// ServerlessInstancePrivateEndpoint is autogenerated from the json schema
type ServerlessInstancePrivateEndpoint struct {
	Endpoints           []ServerlessInstancePrivateEndpointEndpoint `json:",omitempty"`
	SrvConnectionString *string                                     `json:",omitempty"`
	Type                *string                                     `json:",omitempty"`
}

// ServerlessInstancePrivateEndpointEndpoint is autogenerated from the json schema
type ServerlessInstancePrivateEndpointEndpoint struct {
	EndpointId   *string `json:",omitempty"`
	ProviderName *string `json:",omitempty"`
	Region       *string `json:",omitempty"`
}

// Link is autogenerated from the json schema
type Link struct {
	Href *string `json:",omitempty"`
	Rel  *string `json:",omitempty"`
}

// ServerlessInstanceProviderSettings is autogenerated from the json schema
type ServerlessInstanceProviderSettings struct {
	BackingProviderName *string `json:",omitempty"`
	ProviderName        *string `json:",omitempty"`
	RegionName          *string `json:",omitempty"`
}
