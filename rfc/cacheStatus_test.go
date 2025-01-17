package rfc

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	souinCtx "github.com/darkweak/souin/context"
	"github.com/darkweak/souin/errors"
)

func TestSetRequestCacheStatus(t *testing.T) {
	h := http.Header{}

	SetRequestCacheStatus(&h, "AHeader", "Souin")
	if h.Get("Cache-Status") != "Souin; fwd=request; detail=AHeader" {
		errors.GenerateError(t, fmt.Sprintf("The Cache-Status must match %s, %s given", "Souin; fwd=request; detail=AHeader", h.Get("Cache-Status")))
	}
	SetRequestCacheStatus(&h, "", "Souin")
	if h.Get("Cache-Status") != "Souin; fwd=request; detail=" {
		errors.GenerateError(t, fmt.Sprintf("The Cache-Status must match %s, %s given", "Souin; fwd=request; detail=", h.Get("Cache-Status")))
	}
	SetRequestCacheStatus(&h, "A very long header with spaces", "Souin")
	if h.Get("Cache-Status") != "Souin; fwd=request; detail=A very long header with spaces" {
		errors.GenerateError(t, fmt.Sprintf("The Cache-Status must match %s, %s given", "Souin; fwd=request; detail=A very long header with spaces", h.Get("Cache-Status")))
	}
}

func TestValidateCacheControl(t *testing.T) {
	rq := httptest.NewRequest(http.MethodGet, "/", nil)
	rq = rq.WithContext(context.WithValue(rq.Context(), souinCtx.CacheName, "Souin"))
	r := http.Response{
		Request: rq,
	}
	r.Header = http.Header{}

	valid := ValidateCacheControl(&r)
	if !valid {
		errors.GenerateError(t, "The Cache-Control should be valid while an empty string is provided")
	}
	h := http.Header{
		"Cache-Control": []string{"stale-if-error;malformed"},
	}
	r.Header = h
	valid = ValidateCacheControl(&r)
	if valid {
		errors.GenerateError(t, "The Cache-Control shouldn't be valid with max-age")
	}
}

func TestSetCacheStatusEventually(t *testing.T) {
	rq := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(rq.Context(), souinCtx.CacheName, "This")
	ctx = context.WithValue(ctx, souinCtx.Key, "My-key")
	ctx = context.WithValue(ctx, souinCtx.DisplayableKey, true)
	rq = rq.WithContext(ctx)
	r := http.Response{
		Request: rq,
	}
	r.Header = http.Header{}

	SetCacheStatusEventually(&r)
	if r.Header.Get("Cache-Status") != "This; hit; ttl=-1; key=My-key" {
		errors.GenerateError(t, fmt.Sprintf("The Cache-Status should be equal to This; hit; ttl=-1; key=My-key, %s given", r.Header.Get("Cache-Status")))
	}

	r.Header = http.Header{"Date": []string{"Invalid"}}
	SetCacheStatusEventually(&r)
	if r.Header.Get("Cache-Status") == "" {
		errors.GenerateError(t, "The Cache-Control shouldn't be empty")
	}
	if r.Header.Get("Cache-Status") != "This; fwd=request; detail=MALFORMED-DATE" {
		errors.GenerateError(t, "The Cache-Control should be equal to MALFORMED-DATE")
	}

	r.Header = http.Header{}
	SetCacheStatusEventually(&r)
	if ti, e := http.ParseTime(r.Header.Get("Date")); r.Header.Get("Date") == "" || e != nil || r.Header.Get("Date") != ti.Format(http.TimeFormat) {
		errors.GenerateError(t, fmt.Sprintf("Date cannot be null when invalid and must match %s, %s given", r.Header.Get("Date"), ti.Format(http.TimeFormat)))
	}
}
