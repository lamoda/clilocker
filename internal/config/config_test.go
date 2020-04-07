package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewSuccess(t *testing.T) {
	cases := []struct {
		args     []string
		expected *Config
	}{
		{
			args: []string{"ls", "-la"},
			expected: &Config{
				LockerDsn: "local://",
				ConfigFile: "",
				Limit:     0,
				CommandId: "ls -la",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--limit=10", "ls", "-la"},
			expected: &Config{
				LockerDsn: "local://",
				ConfigFile: "",
				Limit:     10,
				CommandId: "ls -la",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--id=test", "ls", "-la"},
			expected: &Config{
				LockerDsn: "local://",
				ConfigFile: "",
				Limit:     0,
				CommandId: "test",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--limit=10", "--id=test", "ls", "-la"},
			expected: &Config{
				LockerDsn: "local://",
				ConfigFile: "",
				Limit:     10,
				CommandId: "test",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--locker_dsn=redis://host:1234", "ls", "-la"},
			expected: &Config{
				LockerDsn: "redis://host:1234",
				ConfigFile: "",
				Limit:     0,
				CommandId: "ls -la",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--config=_test/.clilocker.test", "ls", "-la"},
			expected: &Config{
				LockerDsn: "redis://test:1234",
				ConfigFile: "_test/.clilocker.test",
				Limit:     120,
				CommandId: "test_command",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
		{
			args: []string{"--config=_test/.clilocker.test", "--id=cli_command", "ls", "-la"},
			expected: &Config{
				LockerDsn: "redis://test:1234",
				ConfigFile: "_test/.clilocker.test",
				Limit:     120,
				CommandId: "cli_command",
				Command:   "ls",
				Args:      []string{"-la"},
			},
		},
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("[%d]_%v", i, testCase.args), func(t *testing.T) {
			result, err := New(testCase.args)

			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(testCase.expected, result) {
				t.Errorf("expected %v\ngot %v", testCase.expected, result)
			}
		})
	}
}

func TestNewFailed(t *testing.T) {
	cases := []struct {
		args []string
	}{
		{
			args: []string{},
		},
		{
			args: []string{"--limit=10"},
		},
		{
			args: []string{"--unknown-flag=10"},
		},
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("[%d]_%v", i, testCase.args), func(t *testing.T) {
			_, err := New(testCase.args)

			if err == nil {
				t.Errorf("expected error")
			}
		})
	}
}
