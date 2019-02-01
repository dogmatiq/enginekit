package message

// Correlation is holds identifiers for a specific message.
type Correlation struct {
	// MessageID is a unique identifier for the message.
	MessageID string

	// CausationID is the ID of the message that was being handled when the message
	// identified by MessageID was produced.
	CausationID string

	// CorrelationID is the ID of the "root" message that entered the application
	// to cause the message identified by MessageID, either directly or indirectly.
	CorrelationID string
}

// NewCorrelation returns a correlation for the "root" message with the given
// ID.
func NewCorrelation(id string) Correlation {
	return Correlation{id, id, id}
}

// MustValidate panics if c is not a valid correlation.
func (c Correlation) MustValidate() {
	if c.MessageID == "" {
		panic("message ID is empty")
	}

	if c.CausationID == "" {
		panic("causation ID is empty")
	}

	if c.CorrelationID == "" {
		panic("correlation ID is empty")
	}
}

// New returns a correlation for a message caused by c.MessageID.
func (c Correlation) New(id string) Correlation {
	if id == c.MessageID {
		panic("id must not be the same as the parent's message ID")
	}

	if id == c.CausationID {
		panic("id must not be the same as the parent's causation ID")
	}

	if id == c.CorrelationID {
		panic("id must not be the same as the parent's correlation ID")
	}

	return Correlation{id, c.MessageID, c.CorrelationID}
}

// IsRoot returns true if this message is the "root" of a message tree.
func (c Correlation) IsRoot() bool {
	c.MustValidate()
	return c.MessageID == c.CausationID
}

// IsCausedBy returns true if c.MessageID is directly caused by p.
func (c Correlation) IsCausedBy(p Correlation) bool {
	c.MustValidate()
	p.MustValidate()
	return c.CausationID == p.MessageID
}

// IsCorrelatedWith returns true if c.MessageID is correlated with p.
func (c Correlation) IsCorrelatedWith(p Correlation) bool {
	c.MustValidate()
	p.MustValidate()
	return c.CorrelationID == p.CorrelationID
}
