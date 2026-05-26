package routes


import (
	"net/http"

	handlerUsers "github.com/gattini0928/Equilibrium/internal/handlers/users"
	"github.com/gattini0928/Equilibrium/internal/services/auth"
)

func ConsultationRoutes(mux *http.ServeMux, h *handlerUsers.UserHandler, secret []byte) {
	// Entrar na consulta
	mux.Handle("POST /consultations/from-agenda/{agenda_id}/start",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleStartConsultation)))
	mux.Handle("GET /consultations/{consultation_id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleConsultation)))
	// Salvar Infos da Consulta 
	mux.Handle("POST /consultations/{consultation_id}/save",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleSaveConsultationInfos)))
	// Vincula Profissional ao Paciente
	mux.Handle("POST /consultations/{consultation_id}/link-professional",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleAddProfessionalToPatient)))
	// Detalhes da consulta
	mux.Handle("GET /consultations/consultation-detail/{consultation_id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleConsultationDetail)))
}