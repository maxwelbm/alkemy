package ex01

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Teste do handler utilizando httptest.NewRecorder
func TestHelloHandler(t *testing.T) {
	// Cria uma requisição HTTP para o endpoint
	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	// Cria o ResponseRecorder para capturar a resposta
	// - Ele captura a resposta gerada pelo handler (status code, headers e body)
	//		para que você possa verificar se a saída é a esperada.
	// - Você não precisa iniciar um servidor HTTP completo para testar um
	//      endpoint, o que torna o teste mais rápido e fácil.
	// - O teste é diretamente sobre o comportamento da função handler,
	//      isolando outros fatores como rede ou servidores externos.
	// - omo o ResponseRecorder armazena a resposta em memória, é simples
	//      verificar o status code, headers e o corpo da resposta.
	recorder := httptest.NewRecorder()

	// Chama o handler com o ResponseRecorder e a requisição
	HelloHandler(recorder, req)

	// Verifica o status code da resposta
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status esperado: %v, recebido: %v", http.StatusOK, status)
	}

	// Verifica o corpo da resposta
	expected := `{"message":"Hello, World!"}` + "\n"
	if recorder.Body.String() != expected {
		t.Errorf("Corpo esperado: %v, recebido: %v", expected, recorder.Body.String())
	}
}
