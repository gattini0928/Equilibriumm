package routes

import (
	"net/http"

	handlerUsers "github.com/gattini0928/Equilibrium/internal/handlers/users"
	"github.com/gattini0928/Equilibrium/internal/services/auth"
	"github.com/gattini0928/Equilibrium/internal/middleware"

)


func UserRoutes(mux *http.ServeMux, h *handlerUsers.UserHandler, secret []byte) {

	mux.Handle("/static/",
	http.StripPrefix("/static/",
		http.FileServer(http.Dir("static")),
		),
	)
	
	mux.Handle("/js/",
	http.StripPrefix("/js/",
		http.FileServer(http.Dir("js")),
		),
	)

	mux.Handle("GET /{$}",
		middleware.AuthMiddleware(
			string(secret),
			http.HandlerFunc(h.HandleHome),
		),
	)

	// AUTH
	mux.Handle("/signup",
		middleware.AuthMiddleware(
			string(secret),
			http.HandlerFunc(h.HandleSignup),
		),
	)
	
	mux.Handle("/login",
		middleware.AuthMiddleware(
			string(secret),
			http.HandlerFunc(h.HandleLogin),
		),
	)

	// COMPLETAR PERFIL (JWT)
	mux.Handle("/therapists/profile",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleCompleteTherapist)))
	mux.Handle("/psychiatrists/profile",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleCompletePsychiatrist)))

	// Perfil
	mux.Handle("GET /me",
	auth.JWTMiddleware(secret, http.HandlerFunc(h.HandlePerfil)))

	// LISTAGEM PÚBLICA
	mux.Handle("GET /therapists",
		middleware.AuthMiddleware(
			string(secret),
			http.HandlerFunc(h.HandleAllTherapists),
		),
	)
	mux.Handle("GET /psychiatrists",
		middleware.AuthMiddleware(
			string(secret),
			http.HandlerFunc(h.HandleAllPsychiatrists),
		),
	)

	// DETALHES(Clique no card)
	mux.Handle("GET /therapists/id/{id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleTherapistDetail)))
	mux.Handle("GET /psychiatrists/id/{id}", 
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandlePsychiatristDetail)))

	mux.Handle("POST /therapists/{therapist_id}/agendas/{agenda_id}/reserve", 
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleReserveTherapistAgenda)))
	mux.Handle("POST /psychiatrists/{psychiatrist_id}/agendas/{agenda_id}/reserve", 
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleReservePsychiatristAgenda)))

	// DETALHE DO PACIENTE (JWT)
	mux.Handle("GET /patients/{id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandlePatientDetail)))
	
	// Manipulação de Agenda
	mux.Handle("POST /me/agenda",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleAddAgenda)))
	// Desmarcar Consulta (Doutor)
	mux.Handle("POST /me/agenda/{agenda_id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleDeleteAgenda)))
	// Desmarcar Consulta (Paciente)
	mux.Handle("POST /me/agenda/patient/{agenda_id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleDeleteAgendaPatient)))
	// Desmarcar Consulta(Terapeuta ou Psiquiatra)
	mux.Handle("POST /me/agenda/professional/{agenda_id}",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleDeleteAgendaProfessional)))

	mux.Handle("POST /me/price",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleUpdatePrice)))

	mux.Handle("GET /logout",
		auth.JWTMiddleware(secret, http.HandlerFunc(h.HandleLogout)))
}
