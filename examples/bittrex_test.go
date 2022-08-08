package main

import (
	"flag"
	"os"
	"testing"
)

func TestFlagsHttp(T *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	cases := []struct {
		Name           string
		Args           []string
		ExpectedExit   int
		ExpectedOutput string
	}{
		{"no flags", []string{}, 0, ""},
	}
	for _, tc := range cases {
		flag.CommandLine = flag.NewFlagSet(tc.Name, flag.ExitOnError)
		os.Args = append([]string{tc.Name}, tc.Args...)
		actualExit := realMainHttp()
		if tc.ExpectedExit != actualExit {
			T.Errorf("Wrong exit code for args: %v, expected: %v, got: %v", tc.Args, tc.ExpectedExit, actualExit)
		}
	}
}
