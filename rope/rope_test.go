package rope

import (
	"testing"
)

func TestRopeInit(t *testing.T) {
	rope := New()
	if rope == nil {
		t.Fail()
	}
}
