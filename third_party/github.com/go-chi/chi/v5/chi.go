package chi

import "net/http"

type Mux struct{ *http.ServeMux }

func NewRouter() *Mux { return &Mux{ServeMux: http.NewServeMux()} }

func (m *Mux) Get(pattern string, handler http.HandlerFunc) {
	m.HandleFunc(pattern, methodHandler(http.MethodGet, handler))
}
func (m *Mux) Post(pattern string, handler http.HandlerFunc) {
	m.HandleFunc(pattern, methodHandler(http.MethodPost, handler))
}

func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
