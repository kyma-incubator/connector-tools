package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	dataField = "MSG"
	hmacField = "HMAC"
)

type HMAC struct {
	Key         string
	NextHandler httphandler.Handler
}

func (h *HMAC) validateHMAC(suppliedHmac string, msg string) (bool, error) {
	mac := hmac.New(sha512.New, []byte(h.Key))
	suppliedHmacBytes, err := hex.DecodeString(suppliedHmac)

	if err != nil {
		log.Printf("Conversion of HMAC %s to hex byte array failed with error: %s", suppliedHmac, err.Error())
		return false, err
	}

	mac.Write([]byte(msg))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(suppliedHmacBytes, expectedMAC), nil
}

func (h *HMAC) HandleRequest(r *http.Request, ctx *httphandler.RequestContext) *httphandler.Response {
	r.ParseForm()
	msg := r.Form.Get(dataField)
	suppliedHmac := r.Form.Get(hmacField)
	validationResult, err := h.validateHMAC(suppliedHmac, msg)
	if err != nil {

		log.WithFields(
			ctx.GetLoggerFields(),
		).Error("validation of hmac failed: ", err.Error())

		return &httphandler.Response{
			ResponseCode: 400,
			IsSuccess:    false,
			Response: httphandler.JsonError{
				Message: fmt.Sprint("validation of hmac failed: ", err.Error()),
			},
		}
	}

	if validationResult {
		return h.NextHandler.HandleRequest(r, ctx)
	} else {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Error("validation of hmac failed, either message is not authentic or key is not aligned")

		return &httphandler.Response{
			ResponseCode: 403,
			IsSuccess:    false,
			Response: httphandler.JsonError{
				Message: fmt.Sprint("validation of hmac failed, either message is not authentic or key is not aligned"),
			},
		}
	}
}
