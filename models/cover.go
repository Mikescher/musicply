package models

type CoverHash string

type CoverData struct {
	Hash     CoverHash `json:"hash"`
	MimeType string    `json:"mimeType"`
	Data     []byte    `json:"-"`
}
