package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCookie(t *testing.T) {
	config := `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`
	globalSessions, err := NewManager("cookie", config)
	if err != nil {
		t.Fatal("init cookie session err", err)
	}
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess := globalSessions.SessionStart(w, r)
	err = sess.Set("username", "zhaocloud")
	if err != nil {
		t.Fatal("set error,", err)
	}
	if username := sess.Get("username"); username != "zhaocloud" {
		t.Fatal("get username error")
	}
	sess.SessionRelease(w)
	if cookiestr := w.Header().Get("Set-Cookie"); cookiestr == "" {
		t.Fatal("setcookie error")
	} else {
		parts := strings.Split(strings.TrimSpace(cookiestr), ";")
		for k, v := range parts {
			nameval := strings.Split(v, "=")
			if k == 0 && nameval[0] != "gosessionid" {
				t.Fatal("error")
			}
		}
	}
}
