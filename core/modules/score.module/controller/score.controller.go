package score_controller

import (
	"net/http"

	score_service "loja-back/core/modules/score.module/service"

	"github.com/dalton02/licor/httpkit"
)

func LeadBoard(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	return score_service.GetScoreBoard()
}
