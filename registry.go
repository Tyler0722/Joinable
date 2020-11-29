package main

import (
	"net/http"
	"context"
)

type Registry struct {
	request *http.Request
	sessions map[string]SessionInfo
}

type SessionInfo struct {
	session *Session
	err error
}

func (s *Registry) Get(store PGStore, name string) (session *Session, err error) {
	if info, ok := s.sessions[name]; ok {
		session, err = info.session, info.err
	} else {
		session, err = store.New(s.request, name)
		session.name = name
		s.sessions[name] = SessionInfo{session: session, err: err}
	}
	session.Store = store
	return
}

func GetRegistry(r *http.Request) *Registry {
	ctx := r.Context()
	registry := ctx.Value(0)
	if registry != nil {
		return registry.(*Registry)
	}
	newRegistry := &Registry{ 
		request: r,
		sessions: make(map[string]SessionInfo),
	}
	*r = *r.WithContext(context.WithValue(ctx, 0, newRegistry))
	return newRegistry
}