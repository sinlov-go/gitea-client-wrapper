package gitea_token_client_test

import (
	"fmt"
	"github.com/sinlov-go/unittest-kit/env_kit"
	"github.com/sinlov-go/unittest-kit/unittest_file_kit"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	keyEnvDebug  = "CI_DEBUG"
	keyEnvCiNum  = "CI_NUMBER"
	keyEnvCiKey  = "CI_KEY"
	keyEnvCiKeys = "CI_KEYS"

	EnvKeyGiteaBaseUrl  = "PLUGIN_GITEA_BASE_URL"
	EnvKeyGiteaApiKey   = "PLUGIN_GITEA_API_KEY"
	EnvKeyGiteaInsecure = "PLUGIN_GITEA_INSECURE"
)

var (
	// testBaseFolderPath
	//  test base dir will auto get by package init()
	testBaseFolderPath = ""
	testGoldenKit      *unittest_file_kit.TestGoldenKit

	envDebug  = false
	envCiNum  = 0
	envCiKey  = ""
	envCiKeys []string

	// mustSetArgsAsEnvList
	mustSetArgsAsEnvList = []string{
		EnvKeyGiteaBaseUrl,
		EnvKeyGiteaApiKey,
	}

	valEnvGiteaBaseUrl  = ""
	valEnvGiteaInsecure = false
	valEnvGiteaApiKey   = ""
)

func init() {
	testBaseFolderPath, _ = getCurrentFolderPath()

	envDebug = env_kit.FetchOsEnvBool(keyEnvDebug, false)
	envCiNum = env_kit.FetchOsEnvInt(keyEnvCiNum, 0)
	envCiKey = env_kit.FetchOsEnvStr(keyEnvCiKey, "")
	envCiKeys = env_kit.FetchOsEnvStringSlice(keyEnvCiKeys)

	testGoldenKit = unittest_file_kit.NewTestGoldenKit(testBaseFolderPath)

	valEnvGiteaBaseUrl = env_kit.FetchOsEnvStr(EnvKeyGiteaBaseUrl, "")
	valEnvGiteaApiKey = env_kit.FetchOsEnvStr(EnvKeyGiteaApiKey, "")
	valEnvGiteaInsecure = env_kit.FetchOsEnvBool(EnvKeyGiteaInsecure, false)
}

// test case basic tools start
// getCurrentFolderPath
//
//	can get run path this golang dir
func getCurrentFolderPath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("can not get current file info")
	}
	return filepath.Dir(file), nil
}

// test case basic tools end

func envMustArgsCheck(t *testing.T) bool {
	for _, item := range mustSetArgsAsEnvList {
		if os.Getenv(item) == "" {
			t.Logf("plasee set env: %s, than run test\nfull need set env %v", item, mustSetArgsAsEnvList)
			return true
		}
	}
	return false
}
