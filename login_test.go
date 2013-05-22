package bingo

import (
	"testing"
)

func TestCheck(t *testing.T) {
	t.SkipNow()
	qq := GetQQClient()
	qq.Id = "@18013148370"
	qq.Password = "ace299792458iioo"
	resp := qq.check()
	t.Log(resp)
}

func TestGetLoginSig(t *testing.T) {
	t.SkipNow()
	qq := GetQQClient()
	login_sig, err := qq.get_login_sig()
	if err != nil {
		t.Log(err)
	}
	t.Log(login_sig)
}

func TestLogin(t *testing.T) {
	qq := GetQQClient()
	qq.Id = "@18013148370"
	qq.Password = "ace299792458iioo"
	qq.login()
}
