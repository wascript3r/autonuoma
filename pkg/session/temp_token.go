package session

import (
	"bytes"
	"encoding/gob"
	"time"
)

type TempToken struct {
	SessionID  string
	Expiration time.Time
}

func (t *TempToken) GobEncode() ([]byte, error) {
	b := &bytes.Buffer{}

	vals := []interface{}{t.SessionID, t.Expiration}
	enc := gob.NewEncoder(b)

	for i := range vals {
		err := enc.Encode(vals[i])
		if err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

func (t *TempToken) GobDecode(data []byte) error {
	b := bytes.NewBuffer(data)

	vals := []interface{}{&t.SessionID, &t.Expiration}
	dec := gob.NewDecoder(b)

	for i := range vals {
		err := dec.Decode(vals[i])
		if err != nil {
			return err
		}
	}

	return nil
}
