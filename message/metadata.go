package message

// MetaData is a container for common message meta-data .
type MetaData struct {
	Correlation Correlation
	Type        Type
	Role        Role
	Direction   Direction
}
