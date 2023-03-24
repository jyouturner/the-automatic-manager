package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//RaiseForStatus check the response status code, raise error if it is not 2xx
func RaiseForStatus(resp *http.Response) error {
	code := resp.StatusCode
	statusOk := code >= 200 && code < 300
	if !statusOk {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http response code %d ", code)
		} else {
			return fmt.Errorf("http response code %d %s", code, string(body))
		}

	}
	return nil
}
