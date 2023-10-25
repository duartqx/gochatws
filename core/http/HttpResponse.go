package http

type HttpResponse struct {
	Status int
	Body   interface{}
	Cookie *Cookie
}
