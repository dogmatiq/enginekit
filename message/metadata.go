package message

import (
	"time"

	"github.com/dogmatiq/dogma"
)

// MetaData is a container for common message meta-data .
type MetaData struct {
	Correlation Correlation
	Type        Type
	Role        Role
	Direction   Direction
}

// TimeoutMetaData is a container for common meta-data about a timeout message.
type TimeoutMetaData struct {
	ProcessName string
	ProcessID   string
	TimeoutTime time.Time
}

// Timeout is a container for timeout messages and their meta-data.
type Timeout struct {
	MetaData TimeoutMetaData
	Message  dogma.Message
}
