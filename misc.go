package wx

type ReadOnlyable interface {
	ReadOnly() bool
	SetReadOnly(b bool)
}
