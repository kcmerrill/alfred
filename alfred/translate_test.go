package alfred

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {
	context := &Context{
		Text: TextConfig{
			Success:     "green",
			SuccessIcon: "checkmark",
			Failure:     "red",
			FailureIcon: "x",
		},
	}
	if translate("{{ .Text.Success }}", context) != "green" {
		t.Fatalf(".Text.Success should be green")
	}
	if translate("{{ .Text.SuccessIcon }}", context) != "checkmark" {
		t.Fatalf(".Text.SuccessIcon should be icon")
	}
	if translate("{{ .Text.Failure }}", context) != "red" {
		t.Fatalf(".Text.Failure should be red")
	}
	assert.Equal(t, translate("{{ .Text.FailureIcon }}", context), "x", "Should be X")
}
