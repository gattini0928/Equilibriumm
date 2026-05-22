package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gattini0928/Equilibrium/internal/middleware"
	"github.com/gattini0928/Equilibrium/internal/models"
	"github.com/gattini0928/Equilibrium/internal/services/auth"
	"github.com/gattini0928/Equilibrium/internal/views"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("Corpo da requisição vazio")
	}
	
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error":err.Error()})
}

func CheckID(pathID string, r *http.Request) (int, error) {
	idStr := r.PathValue(pathID)
	if idStr == "" {
		return 0, errors.New("id é obrigatório")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id inválido")
	}

	return id, nil
}

func CheckJWT(w http.ResponseWriter, ctx context.Context) (int, bool) {
	val := ctx.Value(auth.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, "não autorizado")
		return 0, false
	}

	return userID, true
}

func RenderStatusPage(
    w http.ResponseWriter,
    r *http.Request,
    err error,
    statusCode int,
) {

    log.Println("Erro:", err)

    data := models.StatusView{
        ViewData: models.ViewData{
            IsAuth: middleware.IsAuthenticated(r),
        },
        StatusCode: statusCode,
    }

    switch statusCode {

    case http.StatusInternalServerError:
        data.StatusMessage = "Ocorreu um erro no servidor"

        w.WriteHeader(statusCode)
        _ = views.Status500Page(data).Render(r.Context(), w)

    case http.StatusNotFound:
        data.StatusMessage = "Página não encontrada"
	
        w.WriteHeader(statusCode)
        _ = views.Status404Page(data).Render(r.Context(), w)
	
	case http.StatusBadRequest:
		data.StatusMessage = "Há um problema com o que você enviou"
		w.WriteHeader(statusCode)
		_ = views.Status400Page(data).Render(r.Context(), w)
	case http.StatusForbidden:
		data.StatusMessage = "Você não tem permissão para realizar essa ação"
		_ = views.Status403Page(data).Render(r.Context(), w)
    }
}
