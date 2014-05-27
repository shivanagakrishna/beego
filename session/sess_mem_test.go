package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMem(t *testing.T) {
	globalSessions, _ := NewManager("memory", `{"cookieName":"gosessionid","gclifetime":10}`)
	go globalSessions.GC()
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	err := sess.Set("username", "zhaocloud")
	if err != nil {
		t.Fatal("set error,", err)
	}
	if username := sess.Get("username"); username != "zhaocloud" {
		t.Fatal("get username error")
	}
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
