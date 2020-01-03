package error

import (
	"fmt"
	"testing"
)

func TestWrapError(t *testing.T) {

	source := fmt.Errorf("This is a source error")

	wrappedErr := WrapError(source, "This is the target error, code: %d", 4711)

	if wrappedErr.Error() != fmt.Sprintf("This is the target error, code: %d", 4711) {
		t.Errorf("Wrapperd error shoud be %q, but was %q.", wrappedErr.Error(),
			fmt.Sprintf("This is the target error, code: %d", 4711))
	}

	if wrappedErr.Inner != source {
		t.Errorf("Inner error is not wrapped correctly.")
	}

}
