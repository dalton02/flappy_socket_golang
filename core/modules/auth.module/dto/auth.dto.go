package auth_dto

// Estrutura do Usuário
type Usuario struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Senha string `json:"senha"`
}

type LoginRequest struct {
	Nome  string `json:"nome" validator:"required"`
	Senha string `json:"senha" validator:"required"`
}
