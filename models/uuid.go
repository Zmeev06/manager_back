package models

import (
	"encoding/hex"
	"encoding/json"
)

type UUID [16]byte // {{{

func (this UUID) MarshalJSON() ([]byte, error) {
	v := hex.EncodeToString(this[:])
	return json.Marshal(v)
}
func (this *UUID) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	bytes, err := hex.DecodeString(v)
	copy(this[:], bytes)
	return err
} // }}}

