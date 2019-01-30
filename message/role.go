package message

// Role is an enumeration of the roles a message can perform within an engine.
type Role string

const (
	// CommandRole is the role for command messages.
	CommandRole Role = "command"

	// EventRole is the role for event messages.
	EventRole Role = "event"

	// TimeoutRole is the role for timeout messages.
	TimeoutRole Role = "timeout"
)

// Roles is a slice of the valid message roles.
var Roles = []Role{
	CommandRole,
	EventRole,
	TimeoutRole,
}

// MustValidate panics if r is not a valid message role.
func (r Role) MustValidate() {
	switch r {
	case CommandRole:
	case EventRole:
	case TimeoutRole:
	default:
		panic("invalid role: " + r)
	}
}

// Is returns true if r is one of the given roles.
func (r Role) Is(roles ...Role) bool {
	r.MustValidate()

	for _, x := range roles {
		x.MustValidate()

		if r == x {
			return true
		}
	}

	return false
}

// MustBe panics if r is not one of the given roles.
func (r Role) MustBe(roles ...Role) {
	if !r.Is(roles...) {
		panic("unexpected role: " + r)
	}
}

// MustNotBe panics if r is one of the given roles.
func (r Role) MustNotBe(roles ...Role) {
	if r.Is(roles...) {
		panic("unexpected role: " + r)
	}
}

func (r Role) String() string {
	return string(r)
}
