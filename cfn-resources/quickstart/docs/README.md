# MongoDB::Atlas::Quickstart

A quick way to get started with MongoDB Atlas on AWS

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "MongoDB::Atlas::Quickstart",
    "Properties" : {
        "<a href="#atlasapikeysecret" title="AtlasApiKeySecret">AtlasApiKeySecret</a>" : <i>String</i>,
        "<a href="#clustername" title="ClusterName">ClusterName</a>" : <i>String</i>,
        "<a href="#region" title="Region">Region</a>" : <i>String</i>,
        "<a href="#clustersize" title="ClusterSize">ClusterSize</a>" : <i>String</i>,
        "<a href="#addsampledata" title="AddSampleData">AddSampleData</a>" : <i>Boolean</i>,
        "<a href="#connectionstring" title="ConnectionString">ConnectionString</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: MongoDB::Atlas::Quickstart
Properties:
    <a href="#atlasapikeysecret" title="AtlasApiKeySecret">AtlasApiKeySecret</a>: <i>String</i>
    <a href="#clustername" title="ClusterName">ClusterName</a>: <i>String</i>
    <a href="#region" title="Region">Region</a>: <i>String</i>
    <a href="#clustersize" title="ClusterSize">ClusterSize</a>: <i>String</i>
    <a href="#addsampledata" title="AddSampleData">AddSampleData</a>: <i>Boolean</i>
    <a href="#connectionstring" title="ConnectionString">ConnectionString</a>: <i>String</i>
</pre>

## Properties

#### AtlasApiKeySecret

An AWS Secrets Manager Secret with your MongoDB Atlas Api Key.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ClusterName

Name of your MongoDB Atlas Cluster.

_Required_: Yes

_Type_: String

_Maximum_: <code>50</code>

_Update requires_: [Replacement](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-replacement)

#### Region

AWS Region.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ClusterSize

_Required_: No

_Type_: String

_Allowed Values_: <code>M0</code> | <code>M5</code> | <code>M10</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### AddSampleData

Would you like to load sample data into your cluster?

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ConnectionString

Connection string for accessing your new cluster, read only.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the ClusterName.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### IAMRole

IAM Role generated for accessing your new cluster, read only.

#### ConnectionStrings

Returns the <code>ConnectionStrings</code> value.

