package ticket

type Validate interface {
	RawRequest(s interface{}) error
}
