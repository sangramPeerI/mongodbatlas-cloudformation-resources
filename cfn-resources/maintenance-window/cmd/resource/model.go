// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	ApiKeys              *ApiKeyDefinition `json:",omitempty"`
	AutoDeferOnceEnabled *bool             `json:",omitempty"`
	DayOfWeek            *int              `json:",omitempty"`
	GroupId              *string           `json:",omitempty"`
	HourOfDay            *int              `json:",omitempty"`
	StartASAP            *bool             `json:",omitempty"`
}

// ApiKeyDefinition is autogenerated from the json schema
type ApiKeyDefinition struct {
	PrivateKey *string `json:",omitempty"`
	PublicKey  *string `json:",omitempty"`
}
