package plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func request(input interface{}, variables map[string]interface{}) (interface{}, error) {
	req, err := http.NewRequest(variables["method"].(string), variables["url"].(string), strings.NewReader(variables["data"].(string)))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}