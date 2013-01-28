package train

import (
	"io"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	content := ReadAsset(r.URL.Path)
	reader := strings.NewReader(content)
	io.Copy(w, reader)
	return
}
