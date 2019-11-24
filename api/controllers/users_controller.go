package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"fullstack/api/auth"
	"fullstack/api/models"
	"fullstack/api/responses"
	"fullstack/api/utils/formaterror"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	//receive values from body request
	body, err := ioutil.ReadAll(r.Body)

	//check for the error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	//define new model where the body is store
	user := models.User{}

	//transfer plaintext body to user model
	err = json.Unmarshal(body, &user)

	//check for the error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//prepare user
	user.Prepare()

	//validate user
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//try to create user
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	//accessing variable to all user model method
	user := models.User{}

	users, err := user.FindAllUsers(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server Server) GetUser(w http.ResponseWriter, r *http.Request) {
	//receive values from url
	url := mux.Vars(r)

	uid, err := strconv.ParseUint(url["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	searchingUser, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, searchingUser)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//get url from request
	url := mux.Vars(r)
	uid, err := strconv.ParseUint(url["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//receive body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//transfer body to user model
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//extract id
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	//prepare and validate data
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//update user
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
