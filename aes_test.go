package ydapp

import (
	"bytes"
	"testing"
)

const (
	key   = "1234567812345678"
	msg   = "la la la demaxiya"
	appId = "test_appId_1"
)

func TestEnSuccess(t *testing.T) {
	enc, err := AesEncrypt([]byte(msg), []byte(key), appId)
	if err != nil {
		t.Error(err)
	}
	raw, err := AesDecrypt(enc, []byte(key))
	if err != nil {
		t.Error(err)
	}
	if raw.AppId != appId {
		t.Error("Appid not equal")
	}
	if !bytes.Equal(raw.Data, []byte(msg)) {
		t.Error("Data not equal")
	}
}

func TestEnError(t *testing.T) {
	_, err := AesEncrypt([]byte(msg), []byte("1234567890"), appId)
	if err == nil {
		t.Error("Should encrypt failed")
	}
}
