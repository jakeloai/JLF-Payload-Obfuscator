package main

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"
)

func TestEncodeRoundTrips(t *testing.T) {
	in := "' OR 1=1 -- -"
	if got, _ := url.QueryUnescape(encodeURLEncode(in)); got != in {
		t.Errorf("url encode round-trip failed: %q", got)
	}
	if got, _ := hex.DecodeString(encodeHex(in)); string(got) != in {
		t.Errorf("hex round-trip failed")
	}
	if got, err := base64.StdEncoding.DecodeString(encodeBase64(in)); err != nil || string(got) != in {
		t.Errorf("base64 round-trip failed: %v", err)
	}
}

func TestMutateSQLKeywordCaseMixedInput(t *testing.T) {
	// Capitalized keyword must also be replaced (regression test).
	in := "1' Union Select null-- "
	out := mutateSQLKeywordCase(in)
	lower := strings.ToLower(out)
	if strings.Contains(lower, "union") || strings.Contains(lower, "select") {
		t.Errorf("keyword survived mutation: %q", out)
	}
}

func TestXORExpressionDecodes(t *testing.T) {
	fn := "system"
	expr := generateXORExpression(fn, 0x5A)
	// Structural check: hex part must decode to fn XOR key.
	start := strings.Index(expr, `"`) + 1
	end := strings.Index(expr[start:], `"`) + start
	raw, err := hex.DecodeString(expr[start:end])
	if err != nil || len(raw) != len(fn) {
		t.Fatalf("bad hex in XOR expression: %v", err)
	}
	for i := range raw {
		if raw[i]^0x5A != fn[i] {
			t.Fatalf("XOR mismatch at %d", i)
		}
	}
}

func TestNOTExpressionDecodes(t *testing.T) {
	fn := "assert"
	expr := generateNOTExpression(fn)
	start := strings.Index(expr, `"`) + 1
	end := strings.Index(expr[start:], `"`) + start
	raw, err := hex.DecodeString(expr[start:end])
	if err != nil || len(raw) != len(fn) {
		t.Fatalf("bad hex in NOT expression: %v", err)
	}
	for i := range raw {
		if ^raw[i] != fn[i] {
			t.Fatalf("NOT mismatch at %d", i)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	in := []string{"a", "b", "a", " ", "c", "b"}
	out := removeDuplicates(in)
	if len(out) != 3 || out[0] != "a" || out[1] != "b" || out[2] != "c" {
		t.Errorf("dedup failed: %v", out)
	}
}

func TestRandomCasePreservesRuneCount(t *testing.T) {
	in := "SELECT ẞ payload"
	out := randomCase(in)
	if len([]rune(in)) != len([]rune(out)) {
		t.Errorf("rune count changed: %q -> %q", in, out)
	}
}

func TestObfuscateModesNonEmpty(t *testing.T) {
	payload := "' OR 1=1 -- -"
	if len(obfuscateSQL(payload)) == 0 || len(obfuscateXSS("<svg onload=alert(1)>")) == 0 ||
		len(obfuscatePHP("system('id');", "latest")) == 0 || len(obfuscateGeneral(payload)) == 0 {
		t.Error("a mutation mode returned zero variants")
	}
}
