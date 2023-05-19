// @Author: zqy
// @File: propagator.go
// @Date: 2023/5/19 13:58
// @Description todo

package cookie

import "net/http"

type Propagator struct {
	cookieName   string
	cookieOption func(cookie *http.Cookie)
}
type PropagatorOption func(propagator *Propagator)

func NewPropagator(options ...PropagatorOption) *Propagator {
	p := &Propagator{
		cookieName: "sessid",
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}
func WithPropagatorCookieName(cookieName string) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieName = cookieName
	}
}
func (p *Propagator) Inject(id string, w http.ResponseWriter) error {
	c := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
		// 这里不设置过期时间，因为判断过期是在session层面，而不是cookie来判断
		// session是否存在的
		//Expires:
	}
	p.cookieOption(c)
	http.SetCookie(w, c)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(w http.ResponseWriter) error {
	c := &http.Cookie{
		Name: p.cookieName,
		// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	return nil
}
