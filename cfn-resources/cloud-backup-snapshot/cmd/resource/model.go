// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	ApiKeys          *ApiKeyDefinition                                    `json:",omitempty"`
	CloudProvider    *string                                              `json:",omitempty"`
	ClusterName      *string                                              `json:",omitempty"`
	InstanceName     *string                                              `json:",omitempty"`
	CreatedAt        *string                                              `json:",omitempty"`
	Description      *string                                              `json:",omitempty"`
	ExpiresAt        *string                                              `json:",omitempty"`
	FrequencyType    *string                                              `json:",omitempty"`
	GroupId          *string                                              `json:",omitempty"`
	Id               *string                                              `json:",omitempty"`
	IncludeCount     *bool                                                `json:",omitempty"`
	ItemsPerPage     *int                                                 `json:",omitempty"`
	Links            []Link                                               `json:",omitempty"`
	MasterKeyUUID    *string                                              `json:",omitempty"`
	Members          []ApiAtlasDiskBackupShardedClusterSnapshotMemberView `json:",omitempty"`
	MongodVersion    *string                                              `json:",omitempty"`
	PageNum          *int                                                 `json:",omitempty"`
	PolicyItems      []string                                             `json:",omitempty"`
	ReplicaSetName   *string                                              `json:",omitempty"`
	Results          []ApiAtlasDiskBackupShardedClusterSnapshotView       `json:",omitempty"`
	RetentionInDays  *int                                                 `json:",omitempty"`
	SnapshotId       *string                                              `json:",omitempty"`
	SnapshotIds      []string                                             `json:",omitempty"`
	SnapshotType     *string                                              `json:",omitempty"`
	Status           *string                                              `json:",omitempty"`
	StorageSizeBytes *string                                              `json:",omitempty"`
	TotalCount       *float64                                             `json:",omitempty"`
	Type             *string                                              `json:",omitempty"`
}

// ApiKeyDefinition is autogenerated from the json schema
type ApiKeyDefinition struct {
	PrivateKey *string `json:",omitempty"`
	PublicKey  *string `json:",omitempty"`
}

// Link is autogenerated from the json schema
type Link struct {
	Href *string `json:",omitempty"`
	Rel  *string `json:",omitempty"`
}

// ApiAtlasDiskBackupShardedClusterSnapshotMemberView is autogenerated from the json schema
type ApiAtlasDiskBackupShardedClusterSnapshotMemberView struct {
	CloudProvider  *string `json:",omitempty"`
	Id             *string `json:",omitempty"`
	ReplicaSetName *string `json:",omitempty"`
}

// ApiAtlasDiskBackupShardedClusterSnapshotView is autogenerated from the json schema
type ApiAtlasDiskBackupShardedClusterSnapshotView struct {
	CreatedAt        *string                                              `json:",omitempty"`
	Description      *string                                              `json:",omitempty"`
	ExpiresAt        *string                                              `json:",omitempty"`
	FrequencyType    *string                                              `json:",omitempty"`
	Id               *string                                              `json:",omitempty"`
	Links            []Link                                               `json:",omitempty"`
	MasterKeyUUID    *string                                              `json:",omitempty"`
	Members          []ApiAtlasDiskBackupShardedClusterSnapshotMemberView `json:",omitempty"`
	MongodVersion    *string                                              `json:",omitempty"`
	PolicyItems      []string                                             `json:",omitempty"`
	SnapshotIds      []string                                             `json:",omitempty"`
	SnapshotType     *string                                              `json:",omitempty"`
	Status           *string                                              `json:",omitempty"`
	StorageSizeBytes *string                                              `json:",omitempty"`
	Type             *string                                              `json:",omitempty"`
}
