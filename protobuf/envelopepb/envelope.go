package envelopepb

import (
	"errors"
	"fmt"
)

// Validate returns an error if x is not well-formed.
//
// Well-formedness means that all compulsory fields are populated and that no
// incompatible fields are populated.
//
// It is intentionally fairly permissive, so that message meta-data can be
// obtained even if the message is unable to be handled.
//
// It does not perform "deep" validation, such as ensuring messages, times, etc
// can be unmarshaled.
func (x *Envelope) Validate() error {
	if x == nil {
		return errors.New("must not be nil")
	}

	header := x.GetHeader()
	body := x.GetBody()

	if err := header.validate(); err != nil {
		return fmt.Errorf("invalid header: %w", err)
	}

	if err := body.validate(header); err != nil {
		return fmt.Errorf("invalid body: %w", err)
	}

	return nil
}

// Validate returns an error if x is not well-formed.
//
// Well-formedness means that all compulsory fields are populated and that no
// incompatible fields are populated.
//
// It is intentionally fairly permissive, so that message meta-data can be
// obtained even if the message is unable to be handled.
//
// It does not perform "deep" validation, such as ensuring messages, times, etc
// can be unmarshaled.
func (x *MultiEnvelope) Validate() error {
	if x == nil {
		return errors.New("must not be nil")
	}

	header := x.GetHeader()

	if err := header.validate(); err != nil {
		return fmt.Errorf("invalid header: %w", err)
	}

	for i, b := range x.GetBodies() {
		if err := b.validate(header); err != nil {
			return fmt.Errorf("invalid body at index %d: %w", i, err)
		}
	}

	return nil
}

// validate returns an error if x is not well-formed.
func (x *Source) validate() error {
	if x == nil {
		return errors.New("must not be nil")
	}

	if x.Site != nil {
		if err := x.GetSite().Validate(); err != nil {
			return fmt.Errorf("invalid site (%s): %w", x.GetSite(), err)
		}
	}

	if err := x.GetApplication().Validate(); err != nil {
		return fmt.Errorf("invalid application (%s): %w", x.GetApplication(), err)
	}

	if id := x.GetHandler(); id != nil {
		if err := id.Validate(); err != nil {
			return fmt.Errorf("invalid handler (%s): %w", id, err)
		}
	} else if x.GetInstanceId() != "" {
		return errors.New("invalid instance ID: must not be specified without a handler")
	}

	return nil
}

// validate returns an error if x is not well-formed.
func (x *Message) validate() error {
	if x == nil {
		return errors.New("must not be nil")
	}

	if err := x.GetTypeId().Validate(); err != nil {
		return fmt.Errorf("invalid type ID (%s): %w", x.GetTypeId(), err)
	}

	if x.GetDescription() == "" {
		return errors.New("invalid description: must not be empty")
	}

	// Note, x.Data may be empty. A specific application's marshaling logic could
	// conceivably have a message with no data where the message type is
	// sufficient information.

	return nil
}

// validate returns an error if x is not well-formed.
func (x *Header) validate() error {
	if x == nil {
		return errors.New("must not be nil")
	}

	if err := x.GetCausationId().Validate(); err != nil {
		return fmt.Errorf("invalid causation ID (%s): %w", x.GetCausationId(), err)
	}

	if err := x.GetCorrelationId().Validate(); err != nil {
		return fmt.Errorf("invalid correlation ID (%s): %w", x.GetCorrelationId(), err)
	}

	if err := x.GetSource().validate(); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	return nil
}

// validate returns an error if x is not well-formed.
func (x *Body) validate(header *Header) error {
	if x == nil {
		return errors.New("must not be nil")
	}

	if err := x.GetMessageId().Validate(); err != nil {
		return fmt.Errorf("invalid message ID (%s): %w", x.GetMessageId(), err)
	}

	if err := x.GetCreatedAt().CheckValid(); err != nil {
		return fmt.Errorf("invalid created-at time: %w", err)
	}

	if x.ScheduledFor != nil {
		if err := x.GetScheduledFor().CheckValid(); err != nil {
			return fmt.Errorf("invalid scheduled-for time: %w", err)
		}

		source := header.GetSource()
		if source.GetHandler() == nil || source.GetInstanceId() == "" {
			return errors.New("invalid scheduled-for time: must not be specified without a source handler and instance ID")
		}
	}

	if err := x.GetMessage().validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	return nil
}
