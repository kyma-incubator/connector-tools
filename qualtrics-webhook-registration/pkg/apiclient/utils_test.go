package apiclient

import (
	"fmt"
	"testing"
)



func Test_RemoveTrailingSlash (t *testing.T) {
	urlSlash := "https://www.kyma-project.io/"
	urlNoSlash := "https://www.kyma-project.io"


	validUrl := fmt.Sprintf("%s/test/temp", removeTrailingSlash(urlSlash))

	if validUrl != "https://www.kyma-project.io/test/temp" {
		t.Errorf("url should be https://www.kyma-project.io/test/temp, but was: %s", validUrl)
	}
	validUrl = fmt.Sprintf("%s/test/temp", removeTrailingSlash(urlNoSlash))

	if validUrl != "https://www.kyma-project.io/test/temp" {
		t.Errorf("url should be https://www.kyma-project.io/test/temp, but was: %s", validUrl)
	}


}
