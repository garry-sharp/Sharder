package crypt

import (
	"reflect"
	"testing"
)

func TestGetWordIndex(t *testing.T) {
	var res []int
	words := parseMnemonic(tests[0].mnemonic)
	for _, word := range words {
		v, _ := GetWordIndex("en", word)
		res = append(res, v)
	}
	reflect.DeepEqual(res, tests[0].mnemonicindex)
	r1, e1 := GetWordIndex("aa", "west")
	r2, e2 := GetWordIndex("en", "psychic")
	if e1 == nil {
		t.Error("Expected error for unsupported language")
	}
	if e2 == nil {
		t.Error("Expected error for word not found")
	}
	if r1 != -1 && r2 != -1 {
		t.Error("Expected -1 for word not found/unsupported language")
	}

}
func TestGetWordList(t *testing.T) {
	en, enErr := GetWordList("en")
	fr, frErr := GetWordList("fr")
	other, otherErr := GetWordList("aa")
	if len(other) != 0 && otherErr == nil {
		t.Error("Expected error for unsupported language")
	}

	if enErr != nil || frErr != nil {
		t.Errorf("Error getting word list")
	}

	if fr[642] != "élève" {
		t.Errorf("expected élève, got %s", fr[642])
	}

	if en[1995] != "west" {
		t.Errorf("expected west, got %s", en[1995])
	}
	if len(en) != 2048 {
		t.Errorf("expected 2048 words, got %d", len(en))
	}
	if len(fr) != 2048 {
		t.Errorf("expected 2048 words, got %d", len(fr))
	}

}

func TestGetSupportedLanguages(t *testing.T) {
	expected := []string{"en", "es", "fr", "pt", "it", "cz"}
	result := GetSupportedLanguages()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}
