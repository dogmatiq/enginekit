package marshaler

// Packet is a container of marshaled data and its related meta-data.
type Packet struct {
	// MediaType is a MIME media-type describing the content and encoding of the
	// binary data.
	//
	// It must have a "type" parameter that contains the portable name of the
	// type used to represent the data.
	MediaType string

	// Data is the marshaled binary data.
	Data []byte
}

// PortableName returns the portable name of the type represented by the data.
//
// It panics if the media-type does not have a value "type" parameter.
func (p Packet) PortableName() string {
	_, n, err := parseMediaType(p.MediaType)
	if err != nil {
		panic(err)
	}
	return n
}
