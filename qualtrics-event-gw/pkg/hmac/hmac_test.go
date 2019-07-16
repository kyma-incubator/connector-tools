package hmac

import "testing"

const (
	messageInput     = "{\"Status\":\"Complete\",\"SurveyID\":\"SV_1AmkKSDmoZ4XlQx\",\"ResponseID\":\"R_0B7IcnXszjEUHLP\",\"CompletedDate\":\"2019-06-25 07:58:25\",\"BrandID\":\"sapdevelopment\"}"
	hmacTargetResult = "c2c585a404c239aa3063b462cc33551d3ed3ed2de0899f58d456fc0ef7ebef6f3c3af44957989bd0592603e41edde9528f7d9ed09ee128bf719aa12f46e5c598"
	key              = "kyma4ever"
)

func Test_validateHMAC(t *testing.T) {

	subject := HMAC{Key: key}

	result, _ := subject.validateHMAC(hmacTargetResult, messageInput)

	if result == false {
		t.Errorf("HMAC validation failed, hmac should be equal but is different")
	}

	result, _ = subject.validateHMAC(hmacTargetResult, messageInput+"abc")

	if result == true {
		t.Errorf("HMAC validation failed, hmac should be different but is equal")
	}

	subject2 := HMAC{Key: key + "x"}

	result, _ = subject2.validateHMAC(hmacTargetResult, messageInput)

	if result == true {
		t.Errorf("HMAC validation failed, hmac should be different but is equal")
	}

}
