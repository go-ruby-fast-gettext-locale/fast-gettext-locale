// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-fast-gettext-locale/fast-gettext-locale authors

package locale

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// bareTag matches a plain "language[-_]region" tag with 2-3 letter subtags,
// mirroring fast_gettext's format_locale regexp.
var bareTag = regexp.MustCompile(`^([a-zA-Z]{2,3})[-_]([a-zA-Z]{2,3})$`)

// Format normalizes a single locale tag exactly the way fast_gettext's
// format_locale does: a bare "ll[-_]RR" tag becomes "ll_RR" (language lower,
// region upper); anything else (a bare language, a weight fragment, an already
// odd string) is returned unchanged.
func Format(tag string) string {
	m := bareTag.FindStringSubmatch(tag)
	if m == nil {
		return tag
	}
	return strings.ToLower(m[1]) + "_" + strings.ToUpper(m[2])
}

// Normalize is a fuller BCP-47 / POSIX normalization. It trims surrounding
// whitespace, drops an optional ".charset" encoding suffix (e.g. ".UTF-8") and
// an "@modifier" suffix (e.g. "@euro"), converts a "-" separator to "_", and
// then applies [Format]. An empty or separator-only input yields "".
func Normalize(tag string) string {
	tag = strings.TrimSpace(tag)
	if i := strings.IndexByte(tag, '@'); i >= 0 {
		tag = tag[:i]
	}
	if i := strings.IndexByte(tag, '.'); i >= 0 {
		tag = tag[:i]
	}
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}
	return Format(strings.ReplaceAll(tag, "-", "_"))
}

// Ranked is one entry parsed from an Accept-Language header: a normalized
// locale tag and its quality weight in the range (0, 1].
type Ranked struct {
	Tag     string
	Quality float64
}

// ParseAcceptLanguage parses an HTTP Accept-Language header into normalized
// locale tags ranked by descending quality, mirroring fast_gettext's
// formatted_sorted_locales. Whitespace is ignored, an absent ";q=" defaults to
// quality 1.0, and a malformed quality is treated as 1.0. Entries with quality
// 0 (RFC 7231 "not acceptable") and empty tags are dropped. Ties keep their
// original left-to-right order (a stable sort), so the result is deterministic.
func ParseAcceptLanguage(header string) []Ranked {
	header = strings.ReplaceAll(header, " ", "")
	header = strings.ReplaceAll(header, "\t", "")
	var out []Ranked
	for _, part := range strings.Split(header, ",") {
		if part == "" {
			continue
		}
		tag := part
		q := 1.0
		if i := strings.Index(part, ";q="); i >= 0 {
			tag = part[:i]
			if v, err := strconv.ParseFloat(part[i+3:], 64); err == nil {
				q = v
			}
		}
		tag = Format(tag)
		if tag == "" || q <= 0 {
			continue
		}
		out = append(out, Ranked{Tag: tag, Quality: q})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Quality > out[j].Quality })
	return out
}

// Negotiate returns the best available locale for the given Accept-Language
// header, mirroring fast_gettext's best_locale_in.
//
// Candidates are considered in descending quality order. For each candidate an
// exact match against available wins; otherwise the language part (the text
// before the first "_") is matched. A nil available slice means "any locale is
// acceptable", so the top-ranked candidate is returned as-is. An empty (but
// non-nil) slice matches nothing. When nothing matches, "" is returned.
func Negotiate(header string, available []string) string {
	for _, r := range ParseAcceptLanguage(header) {
		if available == nil {
			return r.Tag
		}
		if contains(available, r.Tag) {
			return r.Tag
		}
		if lang := language(r.Tag); lang != r.Tag && contains(available, lang) {
			return lang
		}
	}
	return ""
}

// language returns the language part of a normalized tag: the text before the
// first "_", or the whole tag when there is no region.
func language(tag string) string {
	if i := strings.IndexByte(tag, '_'); i >= 0 {
		return tag[:i]
	}
	return tag
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
