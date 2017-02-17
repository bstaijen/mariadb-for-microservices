package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bstaijen/mariadb-for-microservices/photo-service/config"
	"github.com/bstaijen/mariadb-for-microservices/photo-service/database"
	"github.com/bstaijen/mariadb-for-microservices/shared/util"
	"github.com/urfave/negroni"

	sharedModels "github.com/bstaijen/mariadb-for-microservices/shared/models"
)

func IPCGetPhotos(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.Body == nil {
			util.SendErrorMessage(w, "bad json")
			return
		}

		type Collection struct {
			Objects []*sharedModels.PhotoRequest `json:"requests"`
		}
		col := &Collection{}
		col.Objects = make([]*sharedModels.PhotoRequest, 0)
		err := util.RequestToJSON(r, &col)
		if err != nil {
			log.Fatal(err)
		}

		photos, err := db.GetPhotos(connection, col.Objects)
		if err != nil {
			util.SendError(w, err)
			return
		}
		type Resp struct {
			Photos []*sharedModels.PhotoResponse `json:"results"`
		}
		util.SendOK(w, &Resp{Photos: photos})
	})
}
