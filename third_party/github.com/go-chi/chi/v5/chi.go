package chi

import (
	"net/http"
	"strings"
)

type route struct {
	method, pattern string
	handler         http.HandlerFunc
}
type Mux struct{ routes []route }

func NewRouter() *Mux { return &Mux{} }
func (m *Mux) Get(pattern string, handler http.HandlerFunc) {
	m.routes = append(m.routes, route{http.MethodGet, pattern, handler})
}
func (m *Mux) Post(pattern string, handler http.HandlerFunc) {
	m.routes = append(m.routes, route{http.MethodPost, pattern, handler})
}
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range m.routes {
		if match(route.pattern, r.URL.Path) {
			if r.Method != route.method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			route.handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
func match(pattern, path string) bool {
	pp := strings.Split(strings.Trim(pattern, "/"), "/")
	sp := strings.Split(strings.Trim(path, "/"), "/")
	if len(pp) != len(sp) {
		return false
	}
	for i := range pp {
		if strings.HasPrefix(pp[i], "{") && strings.HasSuffix(pp[i], "}") {
			continue
		}
		if pp[i] != sp[i] {
			return false
		}
	}
	return true
}
