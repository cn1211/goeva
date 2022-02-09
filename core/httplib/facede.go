package httplib

import (
	"bytes"
	"io/ioutil"
)

func Get(url string) []byte {
	res, err := defaultClient.Get(url)
	exp.CheckErr(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	exp.CheckErr(err)

	return body
}

func Post(url string, postData []byte) ([]byte, error) {
	res, err := defaultClient.Post(url, "application/json", bytes.NewReader(postData))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
