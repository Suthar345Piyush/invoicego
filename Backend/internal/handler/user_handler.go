// user handler function

package handler

import (
	"errors"
	"net/http"

	"github.com/Suthar345Piyush/invoicego/internal/middleware"
	"github.com/Suthar345Piyush/invoicego/internal/service"
	"github.com/Suthar345Piyush/invoicego/internal/util"
)

type UserHandler struct {
	userService *service.UserService
}

// function to making user handler

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// function to get the user

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	//getting user

	user, err := h.userService.GetUserByID(claims.UserID)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)

		return
	}

	// if user get's approval , we sent status ok and user data (id)

	util.WriteSuccess(w, http.StatusOK, user, "User retrieved successfully")

}
