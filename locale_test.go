// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-fast-gettext-locale/fast-gettext-locale authors

package locale

import (
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := map[string]string{
		"de-de": "de_DE",
		"DE_ch": "de_CH",
		"pt-BR": "pt_BR",
		"de":    "de",  // bare language: unchanged
		"0.9":   "0.9", // weight fragment: unchanged
		"":      "",    // empty: unchanged
		"en-US": "en_US",
	}
	for in, want := range cases {
		if got := Format(in); got != want {
			t.Errorf("Format(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"en":                   "en",
		"  de-DE.UTF-8@euro  ": "de_DE",
		"fr_FR.UTF-8":          "fr_FR",
		"sr@latin":             "sr",
		"@euro":                "", // nothing left after stripping modifier
		".UTF-8":               "", // nothing left after stripping charset
		"   ":                  "", // whitespace only
		"pt-br":                "pt_BR",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	got := ParseAcceptLanguage("de-de,de;q=0.9,en;q=0.8")
	want := []Ranked{{"de_DE", 1.0}, {"de", 0.9}, {"en", 0.8}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseAcceptLanguage = %#v, want %#v", got, want)
	}

	// Empty header yields no entries (nil).
	if got := ParseAcceptLanguage(""); got != nil {
		t.Errorf("empty header = %#v, want nil", got)
	}

	// Trailing/double comma produces empty parts that are skipped.
	if got := ParseAcceptLanguage("de,"); len(got) != 1 || got[0].Tag != "de" {
		t.Errorf("trailing comma = %#v", got)
	}

	// Malformed quality falls back to 1.0.
	if got := ParseAcceptLanguage("de;q=notanumber"); len(got) != 1 || got[0].Quality != 1.0 {
		t.Errorf("malformed q = %#v", got)
	}

	// q=0 is dropped (not acceptable).
	if got := ParseAcceptLanguage("en;q=0"); got != nil {
		t.Errorf("q=0 = %#v, want nil", got)
	}

	// Empty tag with a weight is dropped.
	if got := ParseAcceptLanguage(";q=0.5"); got != nil {
		t.Errorf("empty tag = %#v, want nil", got)
	}

	// Equal weights keep left-to-right order (stable sort), and whitespace is
	// ignored.
	got2 := ParseAcceptLanguage(" fr , de ")
	if len(got2) != 2 || got2[0].Tag != "fr" || got2[1].Tag != "de" {
		t.Errorf("stable tie order = %#v", got2)
	}
}

func TestNegotiate(t *testing.T) {
	// nil available => accept anything => top-ranked candidate.
	if got := Negotiate("de-de,en", nil); got != "de_DE" {
		t.Errorf("nil available = %q, want de_DE", got)
	}
	// Exact match.
	if got := Negotiate("de-de", []string{"de_DE"}); got != "de_DE" {
		t.Errorf("exact = %q, want de_DE", got)
	}
	// Language-part fallback.
	if got := Negotiate("de-de", []string{"de"}); got != "de" {
		t.Errorf("language fallback = %q, want de", got)
	}
	// Lower-quality exact match wins over an unavailable higher one.
	if got := Negotiate("fr;q=0.9,en;q=0.8", []string{"en"}); got != "en" {
		t.Errorf("ranked pick = %q, want en", got)
	}
	// No region and no match: language(tag)==tag, nothing matches.
	if got := Negotiate("fr", []string{"de"}); got != "" {
		t.Errorf("no match = %q, want empty", got)
	}
	// Empty (non-nil) available matches nothing.
	if got := Negotiate("de", []string{}); got != "" {
		t.Errorf("empty available = %q, want empty", got)
	}
}
