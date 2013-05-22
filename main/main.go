package main

import (
	"bingo"
)

func main() {
	qq := bingo.GetQQClient()
	qq.Id = "@18013148370"
	qq.Password = "ace299792458iioo"
	qq.Login()
}
