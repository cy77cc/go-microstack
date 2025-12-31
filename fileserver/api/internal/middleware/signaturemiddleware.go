package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

type SignatureConfig struct {
	Secret  string
	SkewSec int64
}

type SignatureMiddleware struct {
	cfg SignatureConfig
}

func NewSignatureMiddleware(cfg SignatureConfig) *SignatureMiddleware {
	return &SignatureMiddleware{cfg: cfg}
}

func (m *SignatureMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uidStr := r.Header.Get("X-User-Id")
		tsStr := r.Header.Get("X-Timestamp")
		sig := r.Header.Get("X-Signature")
		if uidStr == "" || tsStr == "" || sig == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !m.checkTimestamp(ts) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !m.verifySignature(r.Method, r.URL.Path, uidStr, tsStr, sig) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "uid", uid)
		next(w, r.WithContext(ctx))
	}
}

func (m *SignatureMiddleware) checkTimestamp(ts int64) bool {
	if m.cfg.SkewSec <= 0 {
		m.cfg.SkewSec = 300
	}
	now := time.Now().Unix()
	diff := now - ts
	if diff < 0 {
		diff = -diff
	}
	return diff <= m.cfg.SkewSec
}

func (m *SignatureMiddleware) verifySignature(method, path, uid, ts, sig string) bool {
	if m.cfg.Secret == "" {
		return false
	}
	canonical := method + "\n" + path + "\n" + uid + "\n" + ts
	want := hmacSha256Hex([]byte(m.cfg.Secret), []byte(canonical))
	return equal(sig, want)
}

func hmacSha256Hex(key, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func equal(a, b string) bool {
	// constant-time compare
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := 0; i < len(a); i++ {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
