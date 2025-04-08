package server

import (
	"encoding/json"
	"fmt"
	"log"
	auth_controller "loja-back/core/modules/auth.module/controller"
	auth_dto "loja-back/core/modules/auth.module/dto"
	score_controller "loja-back/core/modules/score.module/controller"
	"loja-back/core/server/shared"
	"net/http"
	"time"

	"github.com/dalton02/licor/licor"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Players struct {
	Coordenadas Coords `json:"coordenadas"`
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Pontuacao   int    `json:"pontuacao"`
	Skin        string `json:"skin"`
	GameTime    int    `json:"gametime"`
	Status      string `json:"status"`
	SocketCon   *websocket.Conn
}

type Coords struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Rotacao float64 `json:"rotacao"`
}

type SocketMessage struct {
	Tipo     string                 `json:"tipo"`
	Conteudo map[string]interface{} `json:"conteudo"` //Interpretar o json como string no cliente
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     websocket.IsWebSocketUpgrade,
}

var players []Players

func MainServer() {

	go broadcast()

	http.HandleFunc("/upgrade", upgrade)
	corts := cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowPrivateNetwork: true,
	})
	licor.SetCors(corts)
	licor.Public[auth_dto.LoginRequest, any]("/auth/login").Post(auth_controller.Login)
	licor.Public[auth_dto.LoginRequest, any]("/auth/cadastro").Post(auth_controller.Cadastro)
	licor.Public[any, any]("/scoreboard").Get(score_controller.LeadBoard)
	licor.Public[any, any]("/score/{id}").Get(score_controller.UserScore)
	licor.Init("4000")

}
func broadcast() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			emitUpdatedList()
		}
	}
}
func upgrade(w http.ResponseWriter, r *http.Request) {
	var player Players
	c, err := upgrader.Upgrade(w, r, nil)
	player.SocketCon = c
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.SetCloseHandler(func(code int, text string) error {
		return onSocketClose(&player)
	})
	for {
		err := socket(c, &player)
		if err != nil {
			break
		}
	}
}

func onSocketClose(player *Players) error {
	fmt.Println("Jogador desconectado")
	var result []Players
	for _, obj := range players {
		if player.ID != obj.ID {
			result = append(result, obj)
		}
	}
	players = result
	fmt.Println("Lista atualizada de jogadores: ")
	fmt.Println(result)
	emitUpdatedList()
	return nil
}

func socket(c *websocket.Conn, player *Players) error {
	var socketMessage SocketMessage
	var socketWritter SocketMessage

	_, message, err := c.ReadMessage()

	json.Unmarshal([]byte(message), &socketMessage)

	if err != nil {
		log.Println("read:", err)
		return err
	}

	if socketMessage.Tipo == "conectar" {
		jsonData, _ := json.Marshal(socketMessage.Conteudo)
		json.Unmarshal(jsonData, player)
		onConnection(*player)
		socketWritter.Tipo = "conectado"
		socketWritter.Conteudo = map[string]interface{}{
			"mensagem": "Conectado com sucesso a nossa memória :D",
		}
		fmt.Println("Novo jogador conectado, lista atualizada: ")
		fmt.Println(players)
		c.WriteJSON(socketWritter)
	}

	if socketMessage.Tipo == "atualizarPlayer" {
		var playerMemoria Players
		jsonData, _ := json.Marshal(socketMessage.Conteudo)
		json.Unmarshal(jsonData, &playerMemoria)
		for i, obj := range players {
			if playerMemoria.ID == obj.ID {
				players[i].Coordenadas = playerMemoria.Coordenadas
				players[i].Skin = playerMemoria.Skin
				players[i].GameTime = playerMemoria.GameTime
				players[i].Status = playerMemoria.Status
				if players[i].Pontuacao < playerMemoria.Pontuacao {
					players[i].Pontuacao = playerMemoria.Pontuacao
					go updateScore(playerMemoria.Pontuacao, playerMemoria.ID)
				}
			}
		}
	}

	return nil
}

func emitUpdatedList() {

	var socketWritter SocketMessage
	socketWritter.Tipo = "listaAtualizada"
	socketWritter.Conteudo = map[string]interface{}{
		"jogadores": players,
		"mensagem":  "Lista atualizada de jogadores na memoria",
	}
	fmt.Println("Transmissão broadcast: ", players)

	for _, conexao := range players {
		conexao.SocketCon.WriteJSON(socketWritter)
	}

}

func onConnection(player Players) {
	playerExists := false
	for _, obj := range players {
		if player.ID == obj.ID {
			playerExists = true
		}
	}
	if !playerExists {
		players = append(players, player)
	}
}

func updateScore(score int, id int) {
	query := "UPDATE Usuario SET pontuacao = $1 WHERE id = $2"

	_, err := shared.DB.Exec(query, score, id)
	if err != nil {
		fmt.Println("Erro ao atualizar pontuação:", err.Error())
		return
	}
}
