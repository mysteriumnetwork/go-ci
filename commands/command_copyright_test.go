package commands

import (
	"strconv"
	"testing"
)

func TestCopyright(t *testing.T) {
	tests := []struct {
		line  string
		match bool
	}{
		{" * Copyright (C) 2021 The \"MysteriumNetwork/node\" Authors.", true},
		{" * Copyright (c) 2021 The \"MysteriumNetwork/node\" Authors.", true},
		{" * Copyright 2021 The \"MysteriumNetwork/node\" Authors.", false},
	}
	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			if match := copyrightRegex.Match([]byte(tt.line)); match != tt.match {
				t.Errorf("copyrightRegex.Match got = %v, wanted %v", match, tt.match)
			}
		})
	}
}
