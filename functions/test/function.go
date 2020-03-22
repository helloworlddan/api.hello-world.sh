package function

import (
	"fmt"
	"net/http"
)

// Test is an HTTP Cloud Function with a request parameter.
func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Dan!")
}
