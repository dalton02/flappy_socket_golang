package auth_controller

import (
	"encoding/json"
	auth_dto "loja-back/core/modules/auth.module/dto"
	auth_service "loja-back/core/modules/auth.module/service"
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

func Login(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	var loginReq auth_dto.LoginRequest
	json.NewDecoder(request.Body).Decode(&loginReq)
	return auth_service.LoginService(loginReq)

}

func Cadastro(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	var cadReq auth_dto.LoginRequest
	json.NewDecoder(request.Body).Decode(&cadReq)
	return auth_service.CadService(cadReq)

}
