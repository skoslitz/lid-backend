package lidlib

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testPage := new(PageFile)
	testPage.Id = "79_E_898-titel-des-beitrages.md"

	os.Exit(m.Run())
}

func TestReadPage(t *testing.T) {
	page := new(Page)

	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
	}
	for _, test := range tests {
		if got := page.Read(test.input); got != test.want {
			t.Errorf("TestReadPage(%q) = %v", test.input, got)
		}
	}
}
