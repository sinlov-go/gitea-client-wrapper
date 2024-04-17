package gitea_token_client_test

import (
	"github.com/sinlov-go/gitea-client-wrapper/gitea_token_client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientWithHttpTimeout(t *testing.T) {
	if envMustArgsCheck(t) {
		return
	}

	// mock NewClientWithHttpTimeout
	type args struct {
		url           string
		accessToken   string
		timeoutSecond uint
		insecure      bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sample",
			args: args{
				url:           valEnvGiteaBaseUrl,
				accessToken:   valEnvGiteaApiKey,
				timeoutSecond: 30,
				insecure:      valEnvGiteaInsecure,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			giteaTokenClient := gitea_token_client.GiteaTokenClient{}

			// do NewClientWithHttpTimeout
			gotErr := giteaTokenClient.NewClientWithHttpTimeout(tc.args.url, tc.args.accessToken, tc.args.timeoutSecond, tc.args.insecure)

			// verify NewClientWithHttpTimeout
			assert.Equal(t, tc.wantErr, gotErr != nil)
			if tc.wantErr {
				return
			}
			t.Logf("success init as url %s", giteaTokenClient.GetBaseUrl())
		})
	}
}
