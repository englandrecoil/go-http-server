package main

import (
	"testing"
)

func TestRemoveProfanity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "Это kerfuffle!", expected: "Это kerfuffle!"},
		{input: "Я люблю sharbert вокруг костра", expected: "Я люблю **** вокруг костра"},
		{input: "Fornax просто вне этого мира", expected: "**** просто вне этого мира"},
		{input: "Без нецензурной лексики", expected: "Без нецензурной лексики"},
	}

	for _, c := range cases {
		result := removeProfanity(c.input)
		if result != c.expected {
			t.Errorf("Для входной строки '%v' ожидалось '%v', но получили '%v'", c.input, c.expected, result)
		}
	}
}
