// handler function for auth part

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/Suthar345Piyush/invoicego/internal/service"
	"github.com/Suthar345Piyush/invoicego/internal/util"
)

// auth handler struct

type AuthHandler struct {
	authService *service.AuthService
}

// new auth handler function , taking authservice and returns back an auth handler

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

//  auth handler function - REGISTER & LOGIN

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest

	// taking request body to decode it

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)

		return
	}

	// incoming response

	resp, err := h.authService.Register(&req)

	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			util.WriteError(w, http.StatusConflict, err)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//at the end if we got correct response

	util.WriteSuccess(w, http.StatusCreated, resp, "User registered successfully")

}

// LOGIN function

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest

	//passing address of the req for decoding the request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// response completion
	//passing request address to the authservice

	resp, err := h.authService.Login(&req)

	if err != nil {

		if err == domain.ErrInvalidCredentials {
			util.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, err)

		return

	}

	// if response returned correct , then sending correct status with it

	util.WriteSuccess(w, http.StatusOK, resp, "Login Successful")

}
