package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestReadSettingListen(t *testing.T) {

	testCases := []struct {
		desc      string
		flags     []string
		expected  string
		wantError bool
	}{
		{
			desc:      "positive,host and port",
			flags:     []string{`-listen=localhost:8080`},
			expected:  "localhost:8080",
			wantError: false,
		},
		{
			desc:      "positive,port only",
			flags:     []string{`-listen=:8080`},
			expected:  ":8080",
			wantError: false,
		},
		{
			desc:      "negative,missing port",
			flags:     []string{`-listen=host`},
			wantError: true,
		},
		{
			desc:      "negative,missing colon",
			flags:     []string{`-listen=8080`},
			wantError: true,
		},
		{
			desc:      "negative,invalid port",
			flags:     []string{`-listen=:abc`},
			wantError: true,
		},
		{
			desc:      "negative,invalid host and port",
			flags:     []string{`-listen="host 8080"`},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			flagSet := flag.NewFlagSet("testing", flag.ExitOnError)
			err := parseFlags(flagSet, tc.flags)
			if tc.wantError {
				if err == nil {
					t.Errorf("%v: error expected, got nil", tc.desc)
				}
				return
			}
			if err != nil {
				t.Errorf("%v: error not expected, got %q", err, tc.desc)
				return
			}

			if listen != tc.expected {
				t.Errorf("%v: expected %v but got %v", tc.desc, tc.expected, listen)
			}
		})
	}
}

func TestReadSettingRedirect(t *testing.T) {

	// test redirect -
	/*
	   domain.com:http://vm1.dev.test.com/test:301
	   domain.com:vm1.dev.test.com:301
	   
	   domain.com:http://dev.test.com/test:301

	   domain.com:test.com:301
	   domain.com:https://test.com:301
	   domain.com:https://test.com
	   domain.com:test.com

	   test defaults (301) for code
	   test defaults for listern
	*/

	testCases := []struct {
		desc      string
		flags     []string
		expected  map[string]redirectEntry
		wantError bool
	}{
		{
			desc:  "positive, domains only",
			flags: []string{"-redirect=domain.com:test.com"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "", url: "test.com", code: 301},
			},
			wantError: false,
		},
		{
			desc:  "positive, domains with scheme",
			flags: []string{"-redirect=domain.com:https://test.com"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "https", url: "test.com", code: 301},
			},
			wantError: false,
		},
		{
			desc:  "positive, domains with scheme and code",
			flags: []string{"-redirect=domain.com:https://test.com:301"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "https", url: "test.com", code: 301},
			},
			wantError: false,
		},
		{
			desc:  "positive, domains with code",
			flags: []string{"-redirect=domain.com:test.com:302"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "", url: "test.com", code: 302},
			},
			wantError: false,
		},
		{
			desc:  "positive, subsubdomain with code",
			flags: []string{"-redirect=domain.com:vm1.dev.test.com:302"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "", url: "vm1.dev.test.com", code: 302},
			},
			wantError: false,
		},
		{
			desc:  "positive, subdomain with schema and code",
			flags: []string{"-redirect=domain.com:http://dev.test.com/test:301"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "http", url: "dev.test.com/test", code: 301},
			},
			wantError: false,
		},
		{
			desc:  "positive, all options",
			flags: []string{"-redirect=domain.com:http://vm1.dev.test.com/test:302"},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "http", url: "vm1.dev.test.com/test", code: 302},
			},
			wantError: false,
		},
		{
			desc:  "negative, random",
			flags: []string{"-redirect=abccde"},
			wantError: true,
		},
		{
			desc:  "positive, multiple entries",
			flags: []string{
				"-redirect=domain.com:http://vm1.dev.test.com/test:301",
				"-redirect=domain2.com:http://vm3.dev.test2.com/test:302",
			},
			expected: map[string]redirectEntry{
				"domain.com": {scheme: "http", url: "vm1.dev.test.com/test", code: 301},
				"domain2.com": {scheme: "http", url: "vm3.dev.test2.com/test", code: 302},
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			flagSet := flag.NewFlagSet("testing", flag.ExitOnError)
			err := parseFlags(flagSet, tc.flags)
			if tc.wantError {
				if err == nil {
					t.Errorf("%v: error expected, got nil", tc.desc)
				}
				return
			}
			if err != nil {
				t.Errorf("%v: error not expected, got %q", err, tc.desc)
				return
			}

			if !reflect.DeepEqual(redirects, tc.expected) {
				t.Errorf("%v: expected %v but got %v", tc.desc, tc.expected, redirects)
			}
		})
	}

}
