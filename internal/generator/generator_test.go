package generator

import (
	"fmt"
	"testing"
)

func TestGenerateNextSequence(t *testing.T) {
	t.Run("generate next sequence", func(t *testing.T) {
		testCases := []struct {
			current string
			want    string
		}{
			{"aa", "ab"},
			{"aa9", "aba"},
			{"Az9", "AAa"},
			{"99", "aaa"},
			{"999", "aaaa"},
		}

		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%s -> %s", testCase.current, testCase.want), func(t *testing.T) {
				got, err := GenerateNextSequence(testCase.current)
				assertNoError(t, err)
				assertSequenceEqual(t, got, testCase.want)
			})
		}
	})

	t.Run("error on char not from dict", func(t *testing.T) {
		_, err := GenerateNextSequence("fmn_9")
		assertError(t, err)
	})
}

func assertSequenceEqual(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expect an error")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
