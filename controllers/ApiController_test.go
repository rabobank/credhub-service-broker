package controllers

import (
	"testing"
)

var sampleCredentials = map[string]any{
	"a.b.c.d": "value",
	"a": map[string]any{
		"b": map[string]any{
			"c.d": "value",
		},
		"b.c": map[string]any{
			"d": "value",
		},
	},
	"a.b": map[string]any{
		"c.d": "value",
	},
}

func TestDeleteKey(t *testing.T) {
	if !deleteKey(sampleCredentials, "a.b.c.d") {
		t.Error("Failed to delete key a.b.c.d (full)")
	} else if _, found := sampleCredentials["a.b.c.d"]; found {
		t.Error("Full key a.b.c.d present after deletion")
	}

	if !deleteKey(sampleCredentials, "a.b.c.d") {
		t.Error("Failed to delete key a.b.c.d (a b c.d)")
	} else if value, found := sampleCredentials["a"]; !found {
		t.Error("Key a missing after deletion (a b c.d)")
	} else if _, found = value.(map[string]any)["b"]; found {
		t.Error("Key a b not cleaned up after deletion (a b c.d)")
	} else if _, found = value.(map[string]any)["b.c"]; !found {
		t.Error("Key a b.c missing after deletion (a b c.d)")
	}

	if !deleteKey(sampleCredentials, "a.b.c.d") {
		t.Error("Failed to delete key a.b.c.d (a b.c d)")
	} else if _, found := sampleCredentials["a"]; found {
		t.Error("Key a not cleaned up after deletion (a b.c d)")
	}

	if !deleteKey(sampleCredentials, "a.b.c.d") {
		t.Error("Failed to delete key a.b.c.d (a.b c.d)")
	} else if _, found := sampleCredentials["a.b"]; found {
		t.Error("Key a not cleaned up after deletion (a b.c d)")
	}

	if deleteKey(sampleCredentials, "a.b.c.d") {
		t.Error("Delete key a.b.c.d succeeds when no entry is present")
	}

	if len(sampleCredentials) != 0 {
		t.Error("Credentials not empty after all deletions")
	}
}
