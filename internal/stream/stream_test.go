package stream

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestStream_Sanity(t *testing.T) {

	input := map[string]interface{}{
		"Deployment": map[string]interface{}{
			"IsMultiregion": false,
			"JobVersion":    1,
			"Status":        "failed",
			"TaskGroups": map[string]interface{}{
				"group": map[string]interface{}{
					"RequireProgressBy": "2020-10-01T20:20:20.158849655-04:00",
				},
			},
		},
	}

	var out DeployMsg
	var md mapstructure.Metadata
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc(),
		),
		Metadata: &md,
		Result:   &out,
	}

	dec, err := mapstructure.NewDecoder(cfg)
	require.NoError(t, err)

	require.NoError(t, dec.Decode(input))
}
