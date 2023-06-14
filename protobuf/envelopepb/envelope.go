package envelopepb

import (
	"errors"
	"fmt"
)

// Validate returns an error if x is not well-formed.
//
// Well-formedness means that all compulsory fields are populated, and that no
// incompatible fields are populated.
//
// It is intentially fairly permissive, so that message meta-data can be
// obtained even if the message is unable to be handled.
//
// It does not perform "deep" validation, such as ensuring messages, times, etc
// can be unmarshaled.
func (x *Envelope) Validate() error {
	if err := x.GetMessageId().Validate(); err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	if err := x.GetCausationId().Validate(); err != nil {
		return fmt.Errorf("invalid causation ID: %w", err)
	}

	if err := x.GetCorrelationId().Validate(); err != nil {
		return fmt.Errorf("invalid correlation ID: %w", err)
	}

	if x.SourceSite != nil {
		if err := x.GetSourceSite().Validate(); err != nil {
			return fmt.Errorf("invalid source site: %w", err)
		}
	}

	if err := x.GetSourceApplication().Validate(); err != nil {
		return fmt.Errorf("invalid source application: %w", err)
	}

	if id := x.GetSourceHandler(); id != nil {
		if err := id.Validate(); err != nil {
			return fmt.Errorf("invalid source handler: %w", err)
		}
	} else {
		if x.GetSourceInstanceId() != "" {
			return errors.New("invalid source instance ID: must not be specified without a source handler")
		}
		if x.GetScheduledFor() != nil {
			return errors.New("invalid scheduled-for time: must not be specified without a source handler and instance ID")
		}
	}

	if err := x.GetCreatedAt().CheckValid(); err != nil {
		return fmt.Errorf("invalid created-at time: %w", err)
	}

	if x.ScheduledFor != nil {
		if err := x.GetScheduledFor().CheckValid(); err != nil {
			return fmt.Errorf("invalid scheduled-for time: %w", err)
		}
	}

	if x.GetDescription() == "" {
		return errors.New("invalid description: must not be empty")
	}

	if x.GetPortableName() == "" {
		return errors.New("invalid portable name: must not be empty")
	}

	if x.GetMediaType() == "" {
		return errors.New("invalid media-type: must not be empty")
	}

	// Note, e.Data *may* be empty. A specific application's marshaler could
	// conceivably have a message with no data where the message type is
	// sufficient information.

	return nil
}
