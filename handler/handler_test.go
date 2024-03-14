package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var data []byte = []byte(`
{
  "token": "b170",
  "trade": {
    "ticker": "TSLA",
    "strategy": {
      "order": {
        "action": "buy",
        "contracts": "1"
      }
    }
  }
}
`)

func Test_clearTokenValue(t *testing.T) {
	clearData, err := clearTokenValue(data)
	assert.NoError(t, err)
	t.Log(string(clearData))
}
