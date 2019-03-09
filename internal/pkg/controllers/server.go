package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
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

// IndexApiHandler godoc
// @Title Index test
// @Summary Site of OPG+2
// @Description test api handler
// @Produce  json
// @Success 200 {object} models.IndexMessage
// @Router / [get]
func IndexApiHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	message, _ := json.Marshal(models.NewIndexMessage("Site of OPG+2"))
	_, err := fmt.Fprintf(w, string(message))
	if err != nil {
		panic(err)
	}
}
