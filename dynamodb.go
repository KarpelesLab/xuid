package xuid

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// dynamodb specific methods, because aws.......
// see:
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/#Marshaler
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/#Unmarshaler
//
// In order to support Marshal/Unmarshal the AWS api requires depending on the whole aws sdk
// This is stupid, especially since this used to support json marshalers... but aws had to do its own thing...

func (x *XUID) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return errors.New("value missing")
	}
	nv, err := Parse(*av.S)
	if err != nil {
		return err
	}
	x.Prefix = nv.Prefix
	x.UUID = nv.UUID
	return nil
}

func (x XUID) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	val := x.String()
	av.S = &val
	return nil
}
