// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// The identifiers for the temporary security credentials that the operation
// returns.
type AssumedRoleUser struct {

	// The ARN of the temporary security credentials that are returned from the AssumeRole
	// action. For more information about ARNs and how to use them in policies, see [IAM Identifiers]in
	// the IAM User Guide.
	//
	// [IAM Identifiers]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
	//
	// This member is required.
	Arn *string

	// A unique identifier that contains the role ID and the role session name of the
	// role that is being assumed. The role ID is generated by Amazon Web Services when
	// the role is created.
	//
	// This member is required.
	AssumedRoleId *string

	noSmithyDocumentSerde
}

// Amazon Web Services credentials for API authentication.
type Credentials struct {

	// The access key ID that identifies the temporary security credentials.
	//
	// This member is required.
	AccessKeyId *string

	// The date on which the current credentials expire.
	//
	// This member is required.
	Expiration *time.Time

	// The secret access key that can be used to sign requests.
	//
	// This member is required.
	SecretAccessKey *string

	// The token that users must pass to the service API to use the temporary
	// credentials.
	//
	// This member is required.
	SessionToken *string

	noSmithyDocumentSerde
}

// Identifiers for the federated user that is associated with the credentials.
type FederatedUser struct {

	// The ARN that specifies the federated user that is associated with the
	// credentials. For more information about ARNs and how to use them in policies,
	// see [IAM Identifiers]in the IAM User Guide.
	//
	// [IAM Identifiers]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
	//
	// This member is required.
	Arn *string

	// The string that identifies the federated user associated with the credentials,
	// similar to the unique ID of an IAM user.
	//
	// This member is required.
	FederatedUserId *string

	noSmithyDocumentSerde
}

// A reference to the IAM managed policy that is passed as a session policy for a
// role session or a federated user session.
type PolicyDescriptorType struct {

	// The Amazon Resource Name (ARN) of the IAM managed policy to use as a session
	// policy for the role. For more information about ARNs, see [Amazon Resource Names (ARNs) and Amazon Web Services Service Namespaces]in the Amazon Web
	// Services General Reference.
	//
	// [Amazon Resource Names (ARNs) and Amazon Web Services Service Namespaces]: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	Arn *string

	noSmithyDocumentSerde
}

// Contains information about the provided context. This includes the signed and
// encrypted trusted context assertion and the context provider ARN from which the
// trusted context assertion was generated.
type ProvidedContext struct {

	// The signed and encrypted trusted context assertion generated by the context
	// provider. The trusted context assertion is signed and encrypted by Amazon Web
	// Services STS.
	ContextAssertion *string

	// The context provider ARN from which the trusted context assertion was generated.
	ProviderArn *string

	noSmithyDocumentSerde
}

// You can pass custom key-value pair attributes when you assume a role or
// federate a user. These are called session tags. You can then use the session
// tags to control access to resources. For more information, see [Tagging Amazon Web Services STS Sessions]in the IAM User
// Guide.
//
// [Tagging Amazon Web Services STS Sessions]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_session-tags.html
type Tag struct {

	// The key for a session tag.
	//
	// You can pass up to 50 session tags. The plain text session tag keys can’t
	// exceed 128 characters. For these and additional limits, see [IAM and STS Character Limits]in the IAM User
	// Guide.
	//
	// [IAM and STS Character Limits]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-limits.html#reference_iam-limits-entity-length
	//
	// This member is required.
	Key *string

	// The value for a session tag.
	//
	// You can pass up to 50 session tags. The plain text session tag values can’t
	// exceed 256 characters. For these and additional limits, see [IAM and STS Character Limits]in the IAM User
	// Guide.
	//
	// [IAM and STS Character Limits]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-limits.html#reference_iam-limits-entity-length
	//
	// This member is required.
	Value *string

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
