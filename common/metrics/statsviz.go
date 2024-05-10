package metrics

import (
	"github.com/arl/statsviz"
	"net/http"
)

func Serve(addr string) error {
	mux := http.NewServeMux()
	if err := statsviz.Register(mux); err != nil {
		return err
	}
	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}
	return nil
}
