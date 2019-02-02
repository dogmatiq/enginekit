package logging

import (
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

const (
	// MessageIDIcon is the icon shown directly before a message ID.
	// It is an "equals sign", indicating that this message "has exactly" the
	// displayed ID.
	MessageIDIcon = "="

	// CausationIDIcon is the icon shown directly before a message causation ID.
	// It is the mathematical "because" symbol, indicating that this message
	// happened "because of" the displayed ID.
	CausationIDIcon = "∵"

	// CorrelationIDIcon is the icon shown directly before a message correlation ID.
	// It is the mathematical "member of set" symbol, indicating that this message
	// belongs to the set of messages that came about because of the displayed ID.
	CorrelationIDIcon = "⋲"

	// InboundIcon is the icon shown to indicate that a message is "inbound" to a handler.
	// It is a downward pointing arrow, as inbound messages could be considered as
	// being "downloaded" from the network or queue.
	InboundIcon = "▼"

	// InboundErrorIcon is a variant of InboundIcon used when there is an error
	// condition. It is an hollow version of the regular inbound icon, indicating
	// that the requirement remains "unfulfilled".
	InboundErrorIcon = "▽"

	// OutboundIcon is the icon shown to indicate that a message is "outbound" from
	// a handler. It is an upward pointing arrow, as outbound messages could be
	// considered as being "uploaded" to the network or queue.
	OutboundIcon = "▲"

	// OutboundErrorIcon is a variant of OutboundIcon used when there is an error
	// condition. It is an hollow version of the regular inbound icon, indicating
	// that the requirement remains "unfulfilled".
	OutboundErrorIcon = "△"

	// RetryIcon is an icon used instead of InboundIcon when a message is being
	// re-attempted. It is an open-circle with an arrow, indicating that the
	// message has "come around again".
	RetryIcon = "↻"

	// ErrorIcon is the icon shown when logging information about an error.
	// It is a heavy cross, indicating a failure.
	ErrorIcon = "✖"

	// AggregateIcon is the icon shown when a log message relates to an aggregate
	// message handler. It is the mathematical "therefore" symbol, representing the
	// decision making as a result of the message.
	AggregateIcon = "∴"

	// ProcessIcon is the icon shown when a log message relates to a process
	// message handler. It is three horizontal lines, representing the step in a
	// process.
	ProcessIcon = "☰"

	// IntegrationIcon is the icon shown when a log message relates to an
	// integration message handler. It is the relational algebra "join" symbol,
	// representing the integration of two systems.
	IntegrationIcon = "⨝"

	// ProjectionIcon is the icon shown when a log message relates to a projection
	// message handler. It is the mathematical "sum" symbol , representing the
	// aggregation of events.
	ProjectionIcon = "Σ"

	// SystemIcon is an icon shown when a log message relates to the internals of
	// the engine. It is a sprocket, representing the inner workings of the
	// machine.
	SystemIcon = "⚙"

	// SeparatorIcon is an icon used to separate strings of unrelated text inside a
	// log message. It is a large bullet, intended to have a large visual impact.
	SeparatorIcon = "●"
)

// DirectionIcon returns the icon to use for the given message direction.
func DirectionIcon(d message.Direction, isError bool) string {
	d.MustValidate()

	if d == message.InboundDirection {
		if isError {
			return InboundErrorIcon
		}

		return InboundIcon
	}

	if isError {
		return OutboundErrorIcon
	}

	return OutboundIcon
}

// HandlerTypeIcon returns the icon to use for the given handler type.
func HandlerTypeIcon(t handler.Type) string {
	t.MustValidate()

	switch t {
	case handler.AggregateType:
		return AggregateIcon
	case handler.ProcessType:
		return ProcessIcon
	case handler.IntegrationType:
		return IntegrationIcon
	default: // handler.ProjectionType:
		return ProjectionIcon
	}
}
