// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	Enabled   *bool   `json:",omitempty"`
	ProjectId *string `json:",omitempty"`
	ApiKeys   *ApiKey `json:",omitempty"`
}

// ApiKey is autogenerated from the json schema
type ApiKey struct {
	PublicKey  *string `json:",omitempty"`
	PrivateKey *string `json:",omitempty"`
}
