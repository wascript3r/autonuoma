package license

type Validate interface {
	RawRequest(s interface{}) error
}
