package controllers

import (
	"fmt"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err := fmt.Fprintf(w, "Site of OPG+2!")
	if err != nil {
		fmt.Println(err)
	}
}
