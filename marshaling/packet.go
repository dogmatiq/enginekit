package marshaling

// Packet is a container of marshaled data and its related meta-data.
type Packet struct {
	// MediaType is a MIME media-type describing the content and encoding of the
	// binary data.
	MediaType string

	// Data is the marshaled binary data.
	Data []byte
}
