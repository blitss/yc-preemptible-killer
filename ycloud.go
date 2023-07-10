package main

import (
	"context"
	"errors"
	"os"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	iamkey "github.com/yandex-cloud/go-sdk/iamkey"
)

type Ycloud struct {
	Sdk   *ycsdk.SDK
}


type YcloudClient interface {
	DeleteNode(string) error
}

func NewYcloudClient() (ycloud YcloudClient, err error) {
	var sdk *ycsdk.SDK;
	var credentials = os.Getenv("YANDEX_APPLICATION_CREDENTIALS")

	if len(credentials) > 0 {
		key, err := iamkey.ReadFromJSONFile(credentials)

		if err != nil {
			return nil, err
		}

		serviceAccount, err := ycsdk.ServiceAccountKey(key)
		if err != nil {
			return nil, err
		}

		sdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
			Credentials: serviceAccount,
		})
	} else {
		sdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
			// Invoking InstanceServiceAccount automatically requests IAM token and uses
			// it to generate SDK authorization credentials
			Credentials: ycsdk.InstanceServiceAccount(),
		})
	}

	if err != nil {
		return nil, err
	}

	ycloud = &Ycloud{
		Sdk: sdk,
	}

	return ycloud, nil
}

func (ycloud *Ycloud) DeleteNode(nodeId string) (err error) {
	result, err := ycloud.Sdk.Compute().Instance().Delete(context.Background(), &compute.DeleteInstanceRequest{
		InstanceId: nodeId,
	})
	if err != nil {
		return err
	}

	if result.GetError() != nil {
		return errors.New(result.GetError().Message)
	}

	return nil
}