// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-fast-gettext-locale/fast-gettext-locale authors

// Package locale is a pure-Go (no cgo, standard-library only) port of the
// locale-handling behaviour of the Ruby fast_gettext gem, i.e. the parts of
// FastGettext::Storage that accept, normalize and negotiate locales.
//
// It provides:
//
//   - [Format] and [Normalize] for turning a single BCP-47 / POSIX locale tag
//     into fast_gettext's canonical "lang_REGION" form. Format reproduces the
//     gem's exact behaviour (rewrite only a bare "ll[-_]RR" tag); Normalize is
//     a fuller cleanup that also strips a ".charset" encoding suffix and an
//     "@modifier".
//   - [ParseAcceptLanguage] which turns an HTTP Accept-Language header into a
//     quality-ranked list of normalized locale tags, mirroring fast_gettext's
//     formatted_sorted_locales.
//   - [Negotiate] which picks the best available locale for a given
//     Accept-Language header, mirroring fast_gettext's best_locale_in: an exact
//     match wins, otherwise the language part (before the region) is matched,
//     and a nil "available" list means "anything is acceptable".
//
// The package has no dependency on any Ruby runtime; the surface is Go-typed so
// a Ruby binding layer (or the companion fast-gettext translation library) can
// build on it.
package locale
