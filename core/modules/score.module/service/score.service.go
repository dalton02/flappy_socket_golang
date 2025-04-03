package score_service

import (
	"loja-back/core/server/shared"

	"github.com/dalton02/licor/httpkit"
)

func GetScoreBoard() httpkit.HttpMessage {
	type retorno struct {
		ID        int    `json:"id"`
		Nome      string `json:"nome"`
		Pontuacao int    `json:"pontuacao"`
	}
	var usuarios []retorno
	query := "SELECT id, nome, pontuacao FROM Usuario ORDER BY pontuacao DESC"
	rows, err := shared.DB.Query(query)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}

	for rows.Next() {
		var usuario retorno
		rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Pontuacao)
		usuarios = append(usuarios, usuario)
	}

	return httpkit.AppSucess("Sucesso na listagem", usuarios)
}
