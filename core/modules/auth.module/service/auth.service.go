package auth_service

import (
	"database/sql"
	auth_dto "loja-back/core/modules/auth.module/dto"
	"loja-back/core/server/shared"

	"github.com/dalton02/licor/httpkit"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(loginReq auth_dto.LoginRequest) httpkit.HttpMessage {
	var usuario auth_dto.Usuario
	query := "SELECT id, nome, senha FROM Usuario WHERE nome = $1"
	err := shared.DB.QueryRow(query, loginReq.Nome).Scan(&usuario.ID, &usuario.Nome, &usuario.Senha)
	if err != nil {
		if err == sql.ErrNoRows {
			return httpkit.AppNotFound("Usuário não encontrado no sistema")
		}
		return httpkit.AppBadRequest(err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(loginReq.Senha))
	if err != nil {
		return httpkit.AppBadRequest("Senha incorreta, tente novamente")
	}

	return httpkit.AppSucess("Sucesso ao logar", map[string]interface{}{
		"id":   usuario.ID,
		"nome": usuario.Nome,
	})
}

func CadService(cadReq auth_dto.LoginRequest) httpkit.HttpMessage {
	var existingUser auth_dto.Usuario
	queryCheck := "SELECT id FROM Usuario WHERE nome = $1"
	err := shared.DB.QueryRow(queryCheck, cadReq.Nome).Scan(&existingUser.ID)
	if err == nil {
		return httpkit.AppBadRequest("Usuário já existe no sistema")
	}
	if err != sql.ErrNoRows {
		return httpkit.AppBadRequest(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cadReq.Senha), bcrypt.DefaultCost)
	if err != nil {
		return httpkit.AppBadRequest("Erro ao encriptar a senha: " + err.Error())
	}

	queryInsert := "INSERT INTO Usuario (nome, senha, pontuacao) VALUES ($1, $2, $3) RETURNING id"
	var newUserID int
	err = shared.DB.QueryRow(queryInsert, cadReq.Nome, string(hashedPassword), 0).Scan(&newUserID)
	if err != nil {
		return httpkit.AppBadRequest("Erro ao cadastrar usuário: " + err.Error())
	}
	return httpkit.AppSucess("Usuário cadastrado com sucesso", map[string]interface{}{
		"id":   newUserID,
		"nome": cadReq.Nome,
	})
}
