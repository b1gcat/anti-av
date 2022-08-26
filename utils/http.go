package utils

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpGet(url, host string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := http.Client{
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Host = host
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	response, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body := make([]byte, 0)
	if response.StatusCode == 200 {
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(response.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body = append(body, buf...)
			}
		default:
			bodyByte, _ := ioutil.ReadAll(response.Body)
			body = bodyByte
		}
	}
	return body, nil
}
