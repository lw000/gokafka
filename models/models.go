package models

import "encoding/json"

type NewActivityNotify struct {
	Uid      string `json:"uid"`
	DeviceId string `json:"device_id"`
	Money    int    `json:"money"`
	IncNum   int    `json:"inc_num"`

	data   []byte
	err    error
	length int
}

func (p NewActivityNotify) Encode() ([]byte, error) {
	if p.data == nil {
		p.data, p.err = json.Marshal(&p)
		p.length = len(p.data)
	}
	return p.data, p.err
}

func (p NewActivityNotify) Length() int {
	if p.data == nil {
		p.data, p.err = json.Marshal(&p)
		p.length = len(p.data)
	}
	return p.length
}
