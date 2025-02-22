# MongoDB::Atlas::CloudBackupSchedule

An example resource schema demonstrating some basic constructs and validation rules.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "MongoDB::Atlas::CloudBackupSchedule",
    "Properties" : {
        "<a href="#id" title="Id">Id</a>" : <i>String</i>,
        "<a href="#autoexportenabled" title="AutoExportEnabled">AutoExportEnabled</a>" : <i>Boolean</i>,
        "<a href="#useorgandgroupnamesinexportprefix" title="UseOrgAndGroupNamesInExportPrefix">UseOrgAndGroupNamesInExportPrefix</a>" : <i>Boolean</i>,
        "<a href="#export" title="Export">Export</a>" : <i><a href="export.md">Export</a></i>,
        "<a href="#copysettings" title="CopySettings">CopySettings</a>" : <i>[ <a href="apiatlasdiskbackupcopysettingview.md">ApiAtlasDiskBackupCopySettingView</a>, ... ]</i>,
        "<a href="#deletecopiedbackups" title="DeleteCopiedBackups">DeleteCopiedBackups</a>" : <i>[ <a href="apideletecopiedbackupsview.md">ApiDeleteCopiedBackupsView</a>, ... ]</i>,
        "<a href="#policies" title="Policies">Policies</a>" : <i>[ <a href="apipolicyview.md">ApiPolicyView</a>, ... ]</i>,
        "<a href="#referencehourofday" title="ReferenceHourOfDay">ReferenceHourOfDay</a>" : <i>Integer</i>,
        "<a href="#referenceminuteofhour" title="ReferenceMinuteOfHour">ReferenceMinuteOfHour</a>" : <i>Integer</i>,
        "<a href="#restorewindowdays" title="RestoreWindowDays">RestoreWindowDays</a>" : <i>Integer</i>,
        "<a href="#updatesnapshots" title="UpdateSnapshots">UpdateSnapshots</a>" : <i>Boolean</i>,
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
        "<a href="#links" title="Links">Links</a>" : <i>[ <a href="link.md">Link</a>, ... ]</i>
    }
}
</pre>

### YAML

<pre>
Type: MongoDB::Atlas::CloudBackupSchedule
Properties:
    <a href="#id" title="Id">Id</a>: <i>String</i>
    <a href="#autoexportenabled" title="AutoExportEnabled">AutoExportEnabled</a>: <i>Boolean</i>
    <a href="#useorgandgroupnamesinexportprefix" title="UseOrgAndGroupNamesInExportPrefix">UseOrgAndGroupNamesInExportPrefix</a>: <i>Boolean</i>
    <a href="#export" title="Export">Export</a>: <i><a href="export.md">Export</a></i>
    <a href="#copysettings" title="CopySettings">CopySettings</a>: <i>
      - <a href="apiatlasdiskbackupcopysettingview.md">ApiAtlasDiskBackupCopySettingView</a></i>
    <a href="#deletecopiedbackups" title="DeleteCopiedBackups">DeleteCopiedBackups</a>: <i>
      - <a href="apideletecopiedbackupsview.md">ApiDeleteCopiedBackupsView</a></i>
    <a href="#policies" title="Policies">Policies</a>: <i>
      - <a href="apipolicyview.md">ApiPolicyView</a></i>
    <a href="#referencehourofday" title="ReferenceHourOfDay">ReferenceHourOfDay</a>: <i>Integer</i>
    <a href="#referenceminuteofhour" title="ReferenceMinuteOfHour">ReferenceMinuteOfHour</a>: <i>Integer</i>
    <a href="#restorewindowdays" title="RestoreWindowDays">RestoreWindowDays</a>: <i>Integer</i>
    <a href="#updatesnapshots" title="UpdateSnapshots">UpdateSnapshots</a>: <i>Boolean</i>
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
    <a href="#links" title="Links">Links</a>: <i>
      - <a href="link.md">Link</a></i>
</pre>

## Properties

#### Id

Unique identifier of the snapshot.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### AutoExportEnabled

Flag that indicates whether automatic export of cloud backup snapshots to the AWS bucket is enabled.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### UseOrgAndGroupNamesInExportPrefix

Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Export

_Required_: No

_Type_: <a href="export.md">Export</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### CopySettings

List that contains a document for each copy setting item in the desired backup policy.

_Required_: No

_Type_: List of <a href="apiatlasdiskbackupcopysettingview.md">ApiAtlasDiskBackupCopySettingView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### DeleteCopiedBackups

List that contains a document for each deleted copy setting whose backup copies you want to delete.

_Required_: No

_Type_: List of <a href="apideletecopiedbackupsview.md">ApiDeleteCopiedBackupsView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Policies

Rules set for this backup schedule.

_Required_: No

_Type_: List of <a href="apipolicyview.md">ApiPolicyView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ReferenceHourOfDay

UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ReferenceMinuteOfHour

UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### RestoreWindowDays

Number of days back in time you can restore to with Continuous Cloud Backup accuracy. Must be a positive, non-zero integer.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### UpdateSnapshots

Flag indicating if updates to retention in the backup policy were applied to snapshots that Atlas took earlier. 

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ApiKeys

_Required_: Yes

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Links

List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.

_Required_: No

_Type_: List of <a href="link.md">Link</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### ClusterId

Unique identifier of the Atlas cluster.

#### NextSnapshot

Timestamp in the number of seconds that have elapsed since the UNIX epoc when Atlas takes the next snapshot.

#### ID

Returns the <code>ID</code> value.

#### ProjectId

The unique identifier of the project for the Atlas cluster.

#### ClusterName

The name of the Atlas cluster that contains the snapshots you want to retrieve.

#### Links

Returns the <code>Links</code> value.

