// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package commands

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func updateReport(repo string) error {
	uri := "https://goreportcard.com/checks"
	payload := strings.NewReader(fmt.Sprintf("repo=%v", url.QueryEscape(repo)))
	req, _ := http.NewRequest("POST", uri, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Minute * 2,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Goreports responded with status: %v", res.StatusCode)
	}
	return nil
}

// GoReport updates the go report for the given repo
func GoReport(repo string) error {
	err := updateReport(repo)
	if err != nil {
		fmt.Println("Report update failure")
		return err
	}
	fmt.Println("Report updated")
	return nil
}
