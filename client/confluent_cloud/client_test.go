package confluent_cloud

import (
	"fmt"
	"go.dfds.cloud/client/confluent_cloud/service_accounts_v2"
	confluent_util "go.dfds.cloud/client/confluent_cloud/util"
	"log"
	"testing"
)

func TestClientAuthenticate(t *testing.T) {
	client := NewClient()
	err := client.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	cloudSession := confluent_util.NewCloudApiKeySession("")

	resp, err := client.ServiceAccountsV2().GetServiceAccounts(cloudSession, service_accounts_v2.GetServiceAccountsRequest{
		PageSize:  "100",
		PageToken: "",
	})
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	for _, sa := range resp {
		fmt.Println(sa.DisplayName)
	}
}
