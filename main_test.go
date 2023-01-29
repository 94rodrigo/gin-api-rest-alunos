package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/94rodrigo/api-go-gin/controllers"
	"github.com/94rodrigo/api-go-gin/database"
	"github.com/94rodrigo/api-go-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupRodasDeTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "12345678901", RG: "123456789"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestVerificaStatusCodeDaSaudacaoComParametro(t *testing.T) {
	r := SetupRodasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/Rodrigo", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Valores deveriam ser iguais")
	mockDaResposta := `{"API diz:":"E aí, Rodrigo! Tudo certo?"}`
	respostaBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, mockDaResposta, string(respostaBody))
}

func TestListandoTodosOsAlunosHnadler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRodasDeTeste()
	r.GET("/alunos", controllers.ExibeTodosOsAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestBuscaAlunoPorCpfHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRodasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCpf)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678901", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestBuscaAlunoPorIdHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRodasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorId)
	pathDaBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", pathDaBusca, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	var alunoMock models.Aluno
	json.Unmarshal(resp.Body.Bytes(), &alunoMock)
	fmt.Println(alunoMock.Nome)
	assert.Equal(t, "Aluno Teste", alunoMock.Nome)
	assert.Equal(t, "12345678901", alunoMock.CPF)
	assert.Equal(t, "123456789", alunoMock.RG)
}

func TestDeletaAlunoHandle(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	r := SetupRodasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	pathDeBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", pathDeBusca, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestEditaAlunoHandle(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRodasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "47123456789", RG: "123456700"}
	valorJson, _ := json.Marshal(aluno)
	pathParaEditar := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", pathParaEditar, bytes.NewBuffer(valorJson))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	var alunoMockAtualizado models.Aluno
	json.Unmarshal(resp.Body.Bytes(), &alunoMockAtualizado)
	assert.Equal(t, "47123456789", alunoMockAtualizado.CPF)
	assert.Equal(t, "123456700", alunoMockAtualizado.RG)
	assert.Equal(t, "Aluno Teste", alunoMockAtualizado.Nome)
}
