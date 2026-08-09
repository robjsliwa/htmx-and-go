package game

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]*State

	seedItems map[string][]Item
}

func NewSessionStore(seedItems map[string][]Item) *SessionStore {
	return &SessionStore{
		sessions:  make(map[string]*State),
		seedItems: seedItems,
	}
}

const SessionCookieName = "adv_session"

func (s *SessionStore) Get(w http.ResponseWriter, r *http.Request) *State {
	if c, err := r.Cookie(SessionCookieName); err == nil && c.Value != "" {
		s.mu.Lock()
		gs := s.sessions[c.Value]
		s.mu.Unlock()
		if gs != nil {
			return gs
		}
	}

	sid := newSessionID()
	st := NewState(s.seedItems)

	s.mu.Lock()
	s.sessions[sid] = st
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return st
}

func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}