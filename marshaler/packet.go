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
