package train

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	content := ReadAsset(r.URL.Path)
	fmt.Fprintf(w, content)
	return
}
