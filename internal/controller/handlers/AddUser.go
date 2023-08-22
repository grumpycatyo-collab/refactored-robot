package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"refactored-robot/internal/package/models"
	"refactored-robot/internal/service"
)

func (h handler) AddUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var user models.User
	json.Unmarshal(body, &user)
	hashedPassword := service.HashPassword(user.Password)
	user.Password = hashedPassword

	if result := h.DB.Create(&user); result.Error != nil {
		fmt.Println(result.Error)
	} else {
		fmt.Println("User Added")
	}

	// Send a 201 created response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}
