package http

import (
	"io"
	"net/http"
)

func PresentReadme(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(
		writer,
		"Please read the README for instructions: https://github.com/SachaKoohgoli/diceroller_server/blob/main/README.md")
}
