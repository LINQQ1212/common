package response

import (
	"net/http"
)

func Res404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("4404"))
}

func R404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func Res500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("5500"))
}

func R500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func ResNull(w http.ResponseWriter) {
	w.Write([]byte("8"))

}

func R302(w http.ResponseWriter, link string) {
	w.Write([]byte(`<html><head><meta charset="utf-8"><meta http-equiv="refresh" content="0; url=` + link + `" /><script type="text/javascript">window.location.href="` + link + `";</script></head></html>`))
}
