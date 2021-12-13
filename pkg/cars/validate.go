package cars

type Validate interface {
	RawRequest(s interface{}) error
}
