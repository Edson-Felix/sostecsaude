package server

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"
	"wwmm/sostecsaude/server/mydb"

	"github.com/dgrijalva/jwt-go"
)

var logTag = "server: "

// Faz a autenticação do usuário administrador e carrega a página de professor
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/html/admin/login.html")
	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(0)

		if err != nil {
			log.Fatal(logTag + err.Error())
		}

		email, senha := r.FormValue("email"), r.FormValue("senha")

		emails := mydb.GetEmails()

		validCredentials := false

		for _, dbEmail := range emails {
			if dbEmail == email {
				validCredentials = true

				break
			}
		}

		if !validCredentials {
			fmt.Fprintf(w, "Email inválido!")

			return
		}

		dbSenha := mydb.GetSenha(email)

		senhaHash := sha256.Sum256([]byte(senha))

		if dbSenha != string(senhaHash[:]) {
			fmt.Fprintf(w, "Senha inválida!")

			return
		}

		perfil := mydb.GetPerfil(email)

		if perfil == "unidade_saude" {
			fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeSaudeLogin+"<&>"+cfg.UnidadeSaudePassword+"<&>"+email)
		} else if perfil == "unidade_manutencao" {
			fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeManutencaoLogin+"<&>"+cfg.UnidadeManutencaoPassword+"<&>"+email)
		} else if perfil == "unidade_transporte" {
			fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeTransporteLogin+"<&>"+cfg.UnidadeTransportePassword+"<&>"+email)
		}
	}
}

//createToken cria Jason Web Token
func createToken(perfil string, user string, senha string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["perfil"] = perfil
	claims["user"] = user
	claims["senha"] = senha
	claims["validade"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, _ := token.SignedString("sostecsaude.Covid19")

	return tokenString
}

// Envia para o administrador a página de disciplinas
func getPageCadastrar(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/cadastrar.html")
}

// Atualiza o password do administrador
func cadastrar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0)

	if err != nil {
		log.Fatal(logTag + err.Error())
	}

	senha, senhaConfirmacao := r.FormValue("senha"), r.FormValue("senha_confirmacao")

	if senha != senhaConfirmacao {
		fmt.Fprintf(w, "As senhas digitadas são diferentes!")

		return
	}

	email := r.FormValue("email")

	emails := mydb.GetEmails()

	for _, dbEmail := range emails {
		if dbEmail == email {
			fmt.Fprintf(w, "Escolha um outro email!")

			return
		}
	}

	perfil := r.FormValue("perfil")

	senhaHash := sha256.Sum256([]byte(senha))

	mydb.Cadastrar(perfil, email, string(senhaHash[:]))

	if perfil == "unidade_saude" {
		fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeSaudeLogin+"<&>"+cfg.UnidadeSaudePassword+"<&>"+email)
	} else if perfil == "unidade_manutencao" {
		fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeManutencaoLogin+"<&>"+cfg.UnidadeManutencaoPassword+"<&>"+email)
	} else if perfil == "unidade_transporte" {
		fmt.Fprintf(w, perfil+"<&>"+cfg.UnidadeTransporteLogin+"<&>"+cfg.UnidadeTransportePassword+"<&>"+email)
	}
}

// Start http and websockets server
func Start() {
	InitConfig()

	mydb.OpenDB()
	mydb.InitTables()

	log.Println("Starting server...")

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/login", login)
	http.HandleFunc("/get_page_cadastrar", getPageCadastrar)
	http.HandleFunc("/cadastrar", cadastrar)

	/*
		Start Server
	*/

	log.Println("Listening on port " + cfg.ServerPort)

	http.ListenAndServe(":"+cfg.ServerPort, nil)
}

//Clean free resources
func Clean() {
	mydb.Close()
}
