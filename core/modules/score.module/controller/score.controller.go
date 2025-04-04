package score_controller

import (
	"net/http"
	"strconv"

	score_service "loja-back/core/modules/score.module/service"

	"github.com/dalton02/licor/httpkit"
)

func LeadBoard(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	return score_service.GetScoreBoard()
}
func UserScore(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	params, _ := httpkit.GetUrlParams(request)
	id, _ := strconv.Atoi(params.Param["id"])
	return score_service.UserScore(id)
}
