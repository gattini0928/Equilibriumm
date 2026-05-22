package users

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
	"log"

	"github.com/gattini0928/Equilibrium/internal/configs"
	"github.com/gattini0928/Equilibrium/internal/models"
	serviceUsers "github.com/gattini0928/Equilibrium/internal/services/users"
	validators "github.com/gattini0928/Equilibrium/internal/services/validators"

	"github.com/gattini0928/Equilibrium/internal/middleware"
	"github.com/gattini0928/Equilibrium/internal/utils"
	"github.com/gattini0928/Equilibrium/internal/views"
)

func (h *UserHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
    isAuth := middleware.IsAuthenticated(r)
	_ = views.IndexPage(isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	form := models.SignupForm{
		Errors: make(map[string]string),
	}

	data := models.SignupView{
		ViewData: models.ViewData{
			IsAuth: middleware.IsAuthenticated(r),
		},
		Form: form,
	}

	cfg := configs.LoadDBConfig()

	switch r.Method {

	case http.MethodGet:
		_ = views.SignupPage(data).Render(r.Context(), w)
		return

	case http.MethodPost:
		name := r.FormValue("name")
		email := r.FormValue("email")
		ageStr := r.FormValue("age")
		cpf := r.FormValue("cpf")

		cpf = strings.ReplaceAll(cpf, ".", "")
		cpf = strings.ReplaceAll(cpf, "-", "")
		cpf = strings.ReplaceAll(cpf, " ", "")

		role := r.FormValue("role")
		password := r.FormValue("password")

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			form.Errors["age"] = "Idade inválida"
		}

		if err := validators.ValidateName(name); err != nil {
			form.Errors["name"] = err.Error()
		}
		if err := validators.ValidateEmail(email); err != nil {
			form.Errors["email"] = err.Error()
		}
		if err := validators.ValidateCpf(cpf); err != nil {
			form.Errors["cpf"] = err.Error()
		}
		if err := validators.ValidatePassword(password); err != nil {
			form.Errors["password"] = err.Error()
		}

		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.SignupPage(data).Render(r.Context(), w)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			form.Errors["image"] = "Erro ao processar imagem"
		}

		var filename string

		file, handler, err := r.FormFile("image")

		if err != nil {
			form.Errors["image"] = "Imagem obrigatória"
		} else {
			defer file.Close()

			filename = fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)

			filepath := "./static/uploads/" + filename
			dst, err := os.Create(filepath)
			if err != nil {
				form.Errors["image"] = "Erro ao salvar imagem"
			} else {
				defer dst.Close()

				_, err = io.Copy(dst, file)
				if err != nil {
					form.Errors["image"] = "Erro ao salvar imagem"
				}
			}
		}

		if filename == "" {
    		form.Errors["image"] = "Imagem inválida"
		}		

		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.SignupPage(data).Render(r.Context(), w)
			return
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: password,
			Age:      age,
			Cpf:      cpf,
			Role:     role,
			Image:    "/static/uploads/" + filename,
		}

		var patient models.Patient
		var therapist models.Therapist
		var psychiatrist models.Psychiatrist

		switch role {
		case "therapist":
		if err := validators.ValidateAge(age, "therapist"); err != nil {
			form.Errors["age"] = err.Error()
			}
			therapist.Specialty = r.FormValue("specialty")
			therapist.Description = r.FormValue("description")

		case "psychiatrist":
			if err := validators.ValidateAge(age, "psychiatrist"); err != nil {
				form.Errors["age"] = err.Error()
			}
			psychiatrist.CRM = r.FormValue("crm")
			psychiatrist.Description = r.FormValue("description")
		}
		
		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.SignupPage(data).Render(r.Context(), w)
			return
		}

		token, err := h.Service.CreateUser(user, patient, therapist, psychiatrist)
		if err != nil {
			if strings.Contains(err.Error(), "users_email_key") {
				form.Errors["email"] = "Email já cadastrado"
			}
			if strings.Contains(err.Error(), "users_cpf_key") {
				form.Errors["cpf"] = "CPF já cadastrado"
			}
			form.General = "Erro ao criar conta"
			data.Form = form
			_ = views.SignupPage(data).Render(r.Context(), w)
			return
		}	

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(cfg.JWTExpirationInSeconds),
			SameSite: http.SameSiteLaxMode,
		})

		switch user.Role {
		case "therapist":
			http.SetCookie(w, &http.Cookie{
				Name:  "flash",
				Value: "Conta criada, complete seu perfil de terapeuta",
				Path:  "/",
				HttpOnly: true,
				MaxAge:   int(cfg.JWTExpirationInSeconds),
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/therapists/profile", http.StatusSeeOther)
			return

		case "psychiatrist":
			http.SetCookie(w, &http.Cookie{
				Name:  "flash",
				Value: "Conta criada, complete seu perfil de psiquiatra",
				Path:  "/",
				SameSite: http.SameSiteLaxMode,
				HttpOnly: true,
				MaxAge:   int(cfg.JWTExpirationInSeconds),
			})
			http.Redirect(w, r, "/psychiatrists/profile", http.StatusSeeOther)
			return

		default:
			http.Redirect(w, r, "/me", http.StatusSeeOther)
			return
		}
	}
}

func (h *UserHandler) HandleCompleteTherapist(w http.ResponseWriter, r *http.Request) {
	form := models.TherapistForm{
		Errors: make(map[string]string),
	}

	var msg string
	
	if cookie, err := r.Cookie("flash"); err == nil {
		msg = cookie.Value

		http.SetCookie(w, &http.Cookie{
			Name:   "flash",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}


	data := models.TherapistView{
		ViewData: models.ViewData{
			IsAuth: middleware.IsAuthenticated(r),
		},
		Form: form,
		Msg: msg,
	}

	switch r.Method {
	case http.MethodGet:
		_ = views.CompleteTherapistInfoPage(data).Render(r.Context(), w)
		return
	case http.MethodPost:
		specialty := r.FormValue("specialty")
		description := r.FormValue("description")

		var price float64
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err)
			form.General = "Erro ao completar informações"
			data.Form = form
			_ = views.CompleteTherapistInfoPage(data).Render(r.Context(), w)
			return
		}

		if err := validators.ValidatePrice(price); err != nil {
			form.Errors["price"] = err.Error()
		}

		if err := validators.ValidateSpecialty(specialty); err != nil {
			form.Errors["specialty"] = err.Error()
		}

		if err := validators.ValidateDescription(description); err != nil {
			form.Errors["description"] = err.Error()
		}

		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.CompleteTherapistInfoPage(data).Render(r.Context(), w)
			return
		}

		userID, ok := utils.CheckJWT(w, r.Context())
		if !ok {
			return
		}

		err = h.Service.CompleteTherapistSignUp(userID, specialty, description, price)
		if err != nil {
			log.Println(err)
			form.General = "Erro ao completar informações"
			data.Form = form
			_ = views.CompleteTherapistInfoPage(data).Render(r.Context(), w)
			return
		}
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleCompletePsychiatrist(w http.ResponseWriter, r *http.Request) {
	form := models.PsychiatristForm{
		Errors: make(map[string]string),
	}

	var msg string

	if cookie, err := r.Cookie("flash"); err == nil {
		msg = cookie.Value

		http.SetCookie(w, &http.Cookie{
			Name:   "flash",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	data := models.PsychiatristView{
		ViewData: models.ViewData{
			IsAuth: middleware.IsAuthenticated(r),
		},
		Form: form,
		Msg: msg,
	}

	switch r.Method {
	case http.MethodGet:
		_ = views.CompletePsychiatristInfoPage(data).Render(r.Context(), w)
		return
	case http.MethodPost:
		crm := r.FormValue("crm")
		description := r.FormValue("description")

		var price float64
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err)
			form.General = "Erro ao completar informações"
			data.Form = form
			_ = views.CompletePsychiatristInfoPage(data).Render(r.Context(), w)
			return
		}

		if err := validators.ValidatePrice(price); err != nil {
			form.Errors["price"] = err.Error()
		}

		if err := validators.ValidateCrm(crm); err != nil {
			form.Errors["crm"] = err.Error()
		}

		if err := validators.ValidateDescription(description); err != nil {
			form.Errors["description"] = err.Error()
		}

		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.CompletePsychiatristInfoPage(data).Render(r.Context(), w)
			return
		}

		userID, ok := utils.CheckJWT(w, r.Context())
		if !ok {
			return
		}

		err = h.Service.CompletePsychiatristSignUp(userID, crm, description, price)
		if err != nil {
			log.Println("erro no service:", err)
			form.General = "Erro ao completar informações"
			data.Form = form
			_ = views.CompletePsychiatristInfoPage(data).Render(r.Context(), w)
			return
		}
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	form := models.LoginForm{
		Errors: make(map[string]string),
	}

	data := models.LoginView{
		ViewData: models.ViewData{
			IsAuth: middleware.IsAuthenticated(r),
		},
		Form: form,
	}

	switch r.Method {

	case http.MethodGet:
		if cookie, err := r.Cookie("flash"); err == nil {
			data.Msg = cookie.Value

			http.SetCookie(w, &http.Cookie{
				Name:   "flash",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
				SameSite: http.SameSiteLaxMode,
			})
		}

		_ = views.LoginPage(data).Render(r.Context(), w)
		return

	case http.MethodPost:
		form.Email = r.FormValue("email")
		form.Password = r.FormValue("password")

		if form.Email == "" {
			form.Errors["email"] = "Email obrigatório"
		}
		if form.Password == "" {
			form.Errors["password"] = "Senha obrigatória"
		}

		if len(form.Errors) > 0 {
			data.Form = form
			_ = views.LoginPage(data).Render(r.Context(), w)
			return
		}

		_, token, err := h.Service.Login(form.Email, form.Password)
		if err != nil {

			if errors.Is(err, serviceUsers.ErrInvalidPassword) ||
				errors.Is(err, serviceUsers.ErrUserNotFound) {

				form.General = "Email ou senha inválidos"
			} else {
				form.General = "Erro interno"
			}

			data.Form = form
			_ = views.LoginPage(data).Render(r.Context(), w)
			return
		}

		cfg := configs.LoadDBConfig()

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(cfg.JWTExpirationInSeconds),
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:  "flash",
			Value: "Login realizado com sucesso",
			Path:  "/",
			HttpOnly: true,
			MaxAge:   int(cfg.JWTExpirationInSeconds),
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/me", http.StatusSeeOther)
		return
	}
}

func (h *UserHandler) HandlePerfil(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	messages := make(map[string]string)

	data := models.PerfilView{
		ViewData: models.ViewData{
			IsAuth: middleware.IsAuthenticated(r),
		},
		Messages: messages,
		Perfil: models.UserPerfil{},
	}
	
	cookie, err := r.Cookie("flash")

	if err == nil {
		messages["cookie"] = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "flash",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
			SameSite: http.SameSiteLaxMode,
		})
	}

	perfil, err := h.Service.Perfil(userID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	switch p := perfil.(type) {
	case models.PatientDashboard:
		if p.Therapist == nil {
			messages["therapist"] = "Você ainda não tem nenhum terapeuta"
		}
		if p.Psychiatrist == nil{
			messages["psychiatrist"] = "Você ainda não tem nenhum psiquiatra"
		}
		if len(p.Consultations) == 0 {
			messages["consultations"] = "Você ainda não tem nenhuma consulta"
		}

		if len(p.AgendasReserved) == 0 {
			messages["reserved_agendas"] = "Você ainda não tem nenhum horário reservado"
		}

		data.Messages = messages
		_ = views.PatientProfilePage(data, p).Render(r.Context(), w)

	case models.DoctorDashboard:
		if len(p.Agendas) == 0 {
			messages["agendas"] = "Você ainda não tem nenhuma agenda"
		}

		if len(p.AgendasReserved) == 0 {
			messages["reserved_agendas"] = "Você ainda não tem nenhum horário reservado"
		}

		if len(p.Patients) == 0 {
			messages["patients"] = "Você ainda não tem nenhum paciente"
		}
		if len(p.Consultations) == 0 {
			messages["consultations"] = "Você ainda não tem nenhuma consulta"
		}

		data.Messages = messages
		_ = views.DoctorProfilePage(data, p).Render(r.Context(), w)
	}
}

func (h *UserHandler) HandleAllTherapists(w http.ResponseWriter, r *http.Request) {
	therapists, err := h.Service.ListAllTherapists()
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	isAuth := middleware.IsAuthenticated(r)
	_ = views.TherapistsPage(therapists, isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandleAllPsychiatrists(w http.ResponseWriter, r *http.Request) {
	psychiatrists, err := h.Service.ListAllPsychiatrists()
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	isAuth := middleware.IsAuthenticated(r)
	_ = views.PsychiatristsPage(psychiatrists, isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandleTherapistDetail(w http.ResponseWriter, r *http.Request) {
	id, err :=  utils.CheckID("id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err ,http.StatusBadRequest)
		return
	}

	therapist, agendas, err := h.Service.TherapistDetail(id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.RenderStatusPage(w, r, err, http.StatusNotFound)
		return
	}

	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	isAuth := middleware.IsAuthenticated(r)
	_ = views.TherapistDetailPage(therapist, agendas, isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandlePsychiatristDetail(w http.ResponseWriter, r *http.Request) {
	id, err := utils.CheckID("id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	psychiatrist, agendas,  err := h.Service.PsychiatristDetail(id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.RenderStatusPage(w, r, err, http.StatusNotFound)
		return
	}

	isAuth := middleware.IsAuthenticated(r)
	_ = views.PsychiatristDetailPage(psychiatrist, agendas, isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandleAddTherapistToPatient(w http.ResponseWriter, r *http.Request) {
	id, err := utils.CheckID("id", r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err)
		return 
	}

	patientID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.TherapistToPatient(patientID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Terapeuta vinculado com sucesso",
	})
}

func (h *UserHandler) HandleAddPsychiatristToPatient(w http.ResponseWriter, r *http.Request) {
	id, err := utils.CheckID("id", r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err)
		return 
	}

	patientID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.PsychiatristToPatient(patientID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
	"message": "Psiquiatra vinculado com sucesso",
	})
}

func (h *UserHandler) HandlePatientTherapistDetail(w http.ResponseWriter, r *http.Request) {
	patientID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	therapist, err := h.Service.PatientTherapistDetail(patientID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, therapist)
}

func (h *UserHandler) HandlePatientPsychiatristDetail(w http.ResponseWriter, r *http.Request) {
	patientID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	psychiatrist, err := h.Service.PatientPsiquiatristDetail(patientID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, psychiatrist)
}

func (h *UserHandler) HandlePatientDetail(w http.ResponseWriter, r *http.Request) {
	patientID, err := utils.CheckID("id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	doctorID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	patient, err := h.Service.PatientDetail(patientID, doctorID)
	if err != nil {
		utils.RenderStatusPage(w, r, err ,http.StatusInternalServerError)
		return 
	}

	isAuth := middleware.IsAuthenticated(r)
	_ = views.PatientDetailPage(patient, isAuth).Render(r.Context(), w)
}

func (h *UserHandler) HandleAddAgenda(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err := r.ParseForm()
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	var day int
	day, err = strconv.Atoi(r.FormValue("day"))
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	var month int
	month, err = strconv.Atoi(r.FormValue("month"))
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}
	hour := r.FormValue("hour")

	_, err = h.Service.AddAgenda(userID, day, month, hour)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError,)
		return
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)

}

func (h *UserHandler) HandleDeleteAgenda(w http.ResponseWriter, r *http.Request) {
	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.RemoveAgenda(userID, agendaID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return 
	}
	
	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleDeleteAgendaPatient(w http.ResponseWriter, r *http.Request) {
	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.RemoveAgendaPatient(userID, agendaID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return 
	}
	
	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleDeleteAgendaProfessional(w http.ResponseWriter, r *http.Request) {
	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.RemoveAgendaProfessional(userID, agendaID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return 
	}
	
	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleUpdatePrice(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err := r.ParseForm()
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	var price float64
	price, err = strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	err = h.Service.UpdatePrice(userID, price)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)
	
}

func (h *UserHandler) HandleReserveTherapistAgenda(w http.ResponseWriter, r *http.Request) {
	therapistID, err := utils.CheckID("therapist_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.ReserveTherapistAgenda(userID, therapistID, agendaID)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.RenderStatusPage(w, r, err, http.StatusForbidden)
			return
		}
		utils.RenderStatusPage(w, r ,err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *UserHandler) HandleReservePsychiatristAgenda(w http.ResponseWriter, r *http.Request) {
	psychiatristID, err := utils.CheckID("psychiatrist_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}
	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return 
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	err = h.Service.ReserveTherapistAgenda(userID, psychiatristID, agendaID)
	if err != nil {
		utils.RenderStatusPage(w, r ,err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/me", http.StatusSeeOther)

}

func (h *UserHandler) HandleConsultationDetail(w http.ResponseWriter, r *http.Request) {
	consultationID, err := utils.CheckID("consultation_id", r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	consultation, err := h.Service.ShowConsultation(userID, consultationID)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, consultation)
}

func (h *UserHandler) HandleConsultation(w http.ResponseWriter, r *http.Request) {
	consultationID, err := utils.CheckID("consultation_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	data, err := h.Service.ShowConsultationRoom(userID, consultationID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	_ = views.ConsultationPage(data).Render(r.Context(), w)
}

func (h *UserHandler) HandleStartConsultation(w http.ResponseWriter, r *http.Request) {
	agendaID, err := utils.CheckID("agenda_id", r)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
		return
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	consultationID, err := h.Service.StartConsultation(userID, agendaID)
	if err != nil {
		utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/consultations/%d", consultationID), http.StatusSeeOther)
}

func (h *UserHandler) HandleSaveConsultationInfos(w http.ResponseWriter, r *http.Request) {
	consultationID, err := utils.CheckID("consultation_id", r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := utils.CheckJWT(w, r.Context())
	if !ok {
		return
	}

	annotation := r.FormValue("annotation")
	if annotation != "" {
		err = h.Service.SaveConsultationAnnotation(userID, consultationID, annotation)
		if err != nil {
			utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	remedyName := r.FormValue("remedy-name")
	remedyDosage := r.FormValue("remedy-dosage")
	remedyQuantityStr := r.FormValue("remedy-quantity")

	if remedyName != "" && remedyDosage != "" && remedyQuantityStr != "" {
		remedyQuantity, err := strconv.Atoi(remedyQuantityStr)
		if err != nil {
			utils.RenderStatusPage(w, r, err, http.StatusBadRequest)
			return
		}

		err = h.Service.SaveConsultationRemedy(userID, consultationID, remedyName, remedyDosage, remedyQuantity)
		if err != nil {
			utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	diagnosis := r.FormValue("diagnosis")
	if diagnosis != "" {
		err = h.Service.SaveConsultationDiagnosis(userID, consultationID, diagnosis)
		if err != nil {
			utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	author := r.FormValue("author")
	title := r.FormValue("title")

	if author != "" && title != "" {
		err = h.Service.SaveConsultationBook(userID, consultationID, author, title)
		if err != nil {
			utils.RenderStatusPage(w, r, err, http.StatusInternalServerError)
			return
		}
	}
}

func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
        Name:   "token",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}