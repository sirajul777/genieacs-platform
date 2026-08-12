package chi

import "net/http"

type Mux struct{ *http.ServeMux }

func NewRouter() *Mux { return &Mux{ServeMux: http.NewServeMux()} }

func (m *Mux) Get(pattern string, handler http.HandlerFunc) { m.HandleFunc(pattern, handler) }
