/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-ci" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
