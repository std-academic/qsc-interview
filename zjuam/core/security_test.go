package core

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	res := GetEncryptedString("e0d3c43d75f316cdec3be4201981662f2be1cc609da4126c6fc1caaab4ebdeeee019b57916113151c3144afad168fd0d2168f7d737f9a8f9af80d7db55899c23", "10001", "testtest")
	ans := "a89804bb826caa5006ae17d59c2da81e05de99bdbe925c2835d72fc1bd145405b88f60b40f029d44e860f7827663e8904230f16f515e0749619cb2e618932fa7"

	if res != ans {
		t.Errorf("got %s, wanted %s", res, ans)
	}
}
