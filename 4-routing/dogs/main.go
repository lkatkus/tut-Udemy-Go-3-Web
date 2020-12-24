package dogs

import "net/http"

type indexHandler struct{}

func (h indexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {}
