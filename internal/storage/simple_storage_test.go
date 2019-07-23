package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

const dbPathTest = "gosh_test.bolt.db"

func TestSimpleIdentifierStorage(t *testing.T) {
	db, _ := bolt.Open(dbPathTest, 0644, nil)
	defer os.Remove(dbPathTest)

	storage := SimpleIdentifierStorage{"", db}

	testCases := []URL{
		{"aaaa", "http://localhost/foo.bar"},
		{"aaab", "http://google.com/"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("add/retrieve %s", testCase.Id), func(t *testing.T) {
			identifier, idErr := storage.AddURL(testCase.Url)
			assertNoError(t, idErr)
			assertStringEqual(t, identifier, testCase.Id)

			url, urlErr := storage.GetURL(testCase.Id)
			assertNoError(t, urlErr)
			assertSequenceEqual(t, url, testCase.Url)
		})
	}
}

func assertStringEqual(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
