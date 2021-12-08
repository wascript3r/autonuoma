package review

type Validate interface {
	RawRequest(s interface{}) error
}
