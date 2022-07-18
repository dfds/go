package util

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	BASE_ENDPOINT = "https://api.confluent.cloud"
)

func SetAuthHeader(req *http.Request, paramSession Session, clientSession Session) error {
	if paramSession != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", paramSession.GetKey()))
	} else {
		if clientSession == nil {
			return errors.New("no session provided with request, unable to complete")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", clientSession.GetKey()))
	}
	return nil
}
