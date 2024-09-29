package marshaler

import (
	"fmt"
	"mime"
)

// formatMediaType returns the media-type as the base media type and the
// message's portable name as the 'type' parameter.
func formatMediaType(base string, portableName string) string {
	return mime.FormatMediaType(
		base,
		map[string]string{"type": portableName},
	)
}

// parseMediaType returns the media-type and the portable type name encoded in
// the packet's MIME media-type.
func parseMediaType(mediaType string) (string, string, error) {
	mt, params, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return "", "", err
	}

	if n, ok := params["type"]; ok {
		return mt, n, nil
	}

	return "", "", errNoTypeParameter
}

var errNoTypeParameter = fmt.Errorf("the media-type does not specify a 'type' parameter")
