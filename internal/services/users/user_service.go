package users

import (
	"errors"

	"github.com/gattini0928/Equilibrium/internal/configs"
	"github.com/gattini0928/Equilibrium/internal/models"
	"github.com/gattini0928/Equilibrium/internal/services/auth"
	"github.com/gattini0928/Equilibrium/internal/services/validators"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTokenFailed = errors.New("token failed")
)

func (s *UserService) validateConsultationAccess(userID int, c models.Consultation, role string) error {
	switch role {

	case "therapist":
		id, err := s.Repo.GetTherapistIDByUserID(userID)
		if err != nil {
			return err
		}
		if c.TherapistID == nil || *c.TherapistID != id {
			return errors.New("forbidden")
		}

	case "psychiatrist":
		id, err := s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return err
		}
		if c.PsychiatristID == nil || *c.PsychiatristID != id {
			return errors.New("forbidden")
		}

	case "patient":
		id, err := s.Repo.GetPatientIDByUserID(userID)
		if err != nil {
			return err
		}
		if c.PatientID != id {
			return errors.New("forbidden")
		}

	default:
		return errors.New("invalid role")
	}

	return nil
}

func (s *UserService) validateAgendaAccess(userID int, agenda models.Agenda, role string) error {
	switch role {
	case "therapist":
		id, err := s.Repo.GetTherapistIDByUserID(userID)
		if err != nil {
			return err
		}
		if agenda.ProfessionalID != id {
			return errors.New("forbidden")
		}

	case "psychiatrist":
		id, err := s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return err
		}
		if agenda.ProfessionalID != id {
			return errors.New("forbidden")
		}

	case "patient":
		id, err := s.Repo.GetPatientIDByUserID(userID)
		if err != nil {
			return err
		}
		if agenda.PatientID != id {
			return errors.New("forbidden")
		}

	default:
		return errors.New("invalid role")
	}

	return nil
}

func (s *UserService) CreateUser(user models.User, patient models.Patient, therapist models.Therapist, psychiatrist models.Psychiatrist) (string, error) {
	var err error

	err = validators.ValidateName(user.Name)
	if err != nil {
		return "", err
	}
	err = validators.ValidateEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = validators.ValidatePassword(user.Password)
	if err != nil {
		return "",err
	}

	if user.Role == "therapist" || user.Role == "psychiatrist" {
		err = validators.ValidateAge(user.Age, user.Role)
		if err != nil {
			return "", err
		}
	}

	err = validators.ValidateCpf(user.Cpf)
	if err != nil {
		return "", err
	}

	hashPassword, err := validators.HashPassword(user.Password)
	if err != nil {
		return "", ErrInvalidPassword
	}

	user.Password = hashPassword

	err = s.Repo.CreateUserWithProfile(&user, &patient, &therapist, &psychiatrist)
	if err != nil {
		return "", err
	}

	cfg := configs.LoadDBConfig()

	token, err := auth.CreateJWT(s.Secret, user.ID, cfg.JWTExpirationInSeconds)
	if err != nil {
		return "", ErrTokenFailed
	}

	return token, nil

}

func (s *UserService) Login(email string, password string) (models.User, string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, "", err
	}

	if !validators.CheckPasswordHash(password, user.Password){
		return models.User{}, "", ErrInvalidPassword
	}

	cfg := configs.LoadDBConfig()

	token, err := auth.CreateJWT(s.Secret, user.ID, cfg.JWTExpirationInSeconds)
	if err != nil {
		return models.User{}, "", ErrTokenFailed
	}

	return user, token, nil
}

// Completar cadastro do therapeuta
func (s *UserService) CompleteTherapistSignUp(userID int, specialty string, description string, price float64) error {

	err := validators.ValidateSpecialty(specialty)
	if err != nil {
		return err
	}

	err = validators.ValidateDescription(description)
	if err != nil {
		return err
	}

	err = validators.ValidatePrice(price)
	if err != nil {
		return err
	}

	return s.Repo.CompleteTherapist(userID, specialty, description, price)
}

// Completar cadastro do psiquiatra
func (s *UserService) CompletePsychiatristSignUp(userID int, crm string, description string, price float64) error {

	err := validators.ValidateCrm(crm)
	if err != nil {
		return err
	}

	err = validators.ValidateDescription(description)
	if err != nil {
		return err
	}

	err = validators.ValidatePrice(price)
	if err != nil {
		return err
	}

	return s.Repo.CompletePsychiatrist(userID, crm, description, price)
}

// Listagem de todos terapeutas
func (s *UserService) ListAllTherapists() ([]models.DoctorWithUser, error) {
	return s.Repo.GetAllTherapists()
}

// Listagem de todos psiquiatras
func (s *UserService) ListAllPsychiatrists() ([]models.DoctorWithUser, error) {
	return s.Repo.GetAllPsychiatrists()
}

// Detalhes do terapeuta
func (s *UserService) TherapistDetail(userID int) (models.DoctorWithUser, []models.Agenda, error) {

	therapist, err := s.Repo.GetTherapistById(userID)
	if err != nil {
		return models.DoctorWithUser{}, nil, err
	}

	agendas, err := s.Repo.GetTherapistAgenda(therapist.ID, "therapist")
	if err != nil {
		return models.DoctorWithUser{}, nil, err
	}

	return therapist, agendas, nil
}

// Detalhes do psiquiatra
func (s *UserService) PsychiatristDetail(userID int) (models.DoctorWithUser, []models.Agenda, error) {

	psychiatrist, err := s.Repo.GetPsychiatristById(userID)
	if err != nil {
		return models.DoctorWithUser{}, nil, err
	}

	agendas, err := s.Repo.GetPsychiatristAgenda(psychiatrist.ID, "psychiatrist")
	if err != nil {
		return models.DoctorWithUser{}, nil, err
	}

	return psychiatrist, agendas, nil
}

func (s *UserService) TherapistToPatient(patientID, therapistID int) error {
	user, err := s.Repo.GetUserByID(patientID)

	if err != nil {
		return err
	}

	if user.Role != "patient" {
		return errors.New("forbidden")
	}
	
	return s.Repo.AddTherapistToPatient(patientID, therapistID)
}

func (s *UserService) PsychiatristToPatient(patientID, psychiatristID int) error {
	user, err := s.Repo.GetUserByID(patientID)
	if err != nil {
		return err
	}

	if user.Role != "patient" {
		return errors.New("forbidden")
	}

	return s.Repo.AddPsychiatristToPatient(patientID, psychiatristID)

}

// Terapeuta ou Psiquiatra vê os detalhes do paciente
func (s *UserService) PatientDetail(patientID, doctorID int) (models.PatientWithUser, error) {
	user, err := s.Repo.GetUserByID(doctorID)
	if err != nil {
		return models.PatientWithUser{}, err
	}

	if user.Role == "therapist" {
		return s.Repo.GetTherapistPatient(patientID)
	}

	if user.Role == "psychiatrist" {
		return s.Repo.GetPsychiatristPatient(patientID)
	}

	return models.PatientWithUser{}, errors.New("forbidden")
}	

// Detalhes do Terapeuta do Paciente
func (s *UserService) PatientTherapistDetail(userID int) (*models.DoctorWithUser, error) {
	return s.Repo.GetPatientTherapist(userID)
}

// Detalhes do Psiquiatra do Paciente
func (s *UserService) PatientPsiquiatristDetail(userID int) (*models.DoctorWithUser, error) {
	return s.Repo.GetPatientPsychiatrist(userID) 
}

// Listar pacientes do terapeuta ou psiquiatra
func (s *UserService) ListMyPatients(userID int) ([]models.PatientWithUser, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	switch user.Role {
	case "therapist":
		return s.Repo.GetTherapistPatients(userID)

	case "psychiatrist":
		return s.Repo.GetPsychiatristPatients(userID)

	default:
		return nil, errors.New("forbidden")
	}
}

// Funções para agenda e preço
func (s *UserService) AddAgenda(userID int, day int, month int, hour string) (models.Agenda, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return models.Agenda{}, err
	}

	var professionalID int
	var professionalRole string

	switch user.Role {
	case "therapist":
		professionalID, err = s.Repo.GetTherapistIDByUserID(userID)
		professionalRole = "therapist"
		if err != nil {
			return models.Agenda{}, err
		}
	case "psychiatrist":
		professionalID, err = s.Repo.GetPsychiatristIDByUserID(userID)
		professionalRole = "psychiatrist"
		if err != nil {
			return models.Agenda{}, err
		}
	default:
		return models.Agenda{}, errors.New("forbidden")
	}

	return s.Repo.InsertAgenda(professionalID, professionalRole, day, month, hour)
}

func (s *UserService) RemoveAgenda(userID int, agendaID int) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	var professionalID int
	var professionalRole string

	switch user.Role {
	case "therapist":
		professionalID, err = s.Repo.GetTherapistIDByUserID(userID)
		professionalRole = "therapist"
	case "psychiatrist":
		professionalID, err = s.Repo.GetPsychiatristIDByUserID(userID)
		professionalRole = "psychiatrist"
	default:
		return errors.New("forbidden")
	}

	if err != nil {
		return err
	}

	return s.Repo.DeleteAgenda(agendaID, professionalID, professionalRole)
}

func (s *UserService) RemoveAgendaPatient(userID int, agendaID int) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Role != "patient" {
		return errors.New("forbidden")
	}

	patientID, err := s.Repo.GetPatientIDByUserID(userID)
	if err != nil {
		return err
	}

	return s.Repo.UnreserveAgendaPatient(agendaID, patientID)
}

func (s *UserService) RemoveAgendaProfessional(userID int, agendaID int) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Role != "therapist" && user.Role != "psychiatrist" {
		return errors.New("forbidden")
	}

	var professionalID int

	if user.Role == "therapist" {
		professionalID, err = s.Repo.GetTherapistIDByUserID(userID)
		if err != nil {
			return err
		}	
	}

	if user.Role == "psychiatrist" {
		professionalID, err = s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return err
		}
	}

	return s.Repo.UnreserveAgendaProfessional(agendaID, professionalID)
}

func (s *UserService) UpdatePrice(userID int, price float64) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	switch user.Role {
	case "therapist":
		return s.Repo.UpdateTherapistPrice(userID, price)

	case "psychiatrist":
		return s.Repo.UpdatePsychiatristPrice(userID, price)

	default:
		return errors.New("invalid role")
	}
}

func (s *UserService) Perfil(userID int) (any, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	var professionalRole string

	switch user.Role {
	case "patient":

		therapist, err := s.Repo.GetPatientTherapist(userID)
		if err != nil {
			return models.DoctorWithUser{}, err
		}

		psychiatrist, err := s.Repo.GetPatientPsychiatrist(userID)
		if err != nil {
			return models.DoctorWithUser{}, err
		}

		patientID, err := s.Repo.GetPatientIDByUserID(userID)
		if err != nil {
			return models.DoctorWithUser{}, err
		}

		consultations, err := s.Repo.GetPatientConsultations(patientID)
		if err != nil {
			return models.DoctorWithUser{}, err
		}

		agendasReserved, err := s.Repo.GetPatientReservedAgendas(patientID)
		if err != nil {
			return nil, err
		}

		perfil, err := s.Repo.GetPatientPerfil(userID)
		if err != nil {
			return models.PatientDashboard{}, err
		}

		return models.PatientDashboard{
			Perfil: perfil,
			Role: "patient",
			Therapist: therapist,
			Psychiatrist: psychiatrist,
			AgendasReserved: agendasReserved,
			Consultations: consultations,
		}, nil

	case "therapist":
		perfil, err := s.Repo.GetTherapistPerfil(userID)
		if err != nil {
			return nil, err
		}

		therapistID, err := s.Repo.GetTherapistIDByUserID(userID)
		if err != nil {
			return nil, err
		}

		professionalRole = "therapist"

		agendas, err := s.Repo.GetTherapistPrivateAgenda(therapistID, professionalRole)
		if err != nil {
			return nil, err
		}

		agendasReserved, err := s.Repo.GetTherapistReservedAgendas(therapistID, professionalRole)
		if err != nil {
			return nil, err
		}

		patients,  err := s.Repo.GetTherapistPatients(therapistID)
		if err != nil {
			return nil, err
		}

		consultations, err := s.Repo.GetTherapistConsultations(therapistID)
		if err != nil {
			return nil, err
		}

		return models.DoctorDashboard{
			Perfil: perfil,
			Role: "therapist",
			Agendas: agendas,
			AgendasReserved: agendasReserved,
			Patients: patients,
			Consultations: consultations,
		}, nil

	case "psychiatrist":
		perfil, err := s.Repo.GetPsychiatristPerfil(userID)
		if err != nil {
			return nil, err
		}

		psychiatristID , err := s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return nil, err
		}

		professionalRole = "psychiatrist"

		agendas, err := s.Repo.GetPsychiatristPrivateAgenda(psychiatristID, professionalRole)
		if err != nil {
			return nil, err
		}

		agendasReserved, err := s.Repo.GetPsychiatristReservedAgendas(psychiatristID, professionalRole)
		if err != nil {
			return nil, err
		}

		patients,  err := s.Repo.GetPsychiatristPatients(psychiatristID)
		if err != nil {
			return nil, err
		}

		consultations, err := s.Repo.GetPsychiatristConsultations(psychiatristID)
		if err != nil {
			return nil, err
		}
		return models.DoctorDashboard{
			Perfil: perfil,
			Role: "psychiatrist",
			Agendas: agendas,
			AgendasReserved: agendasReserved,
			Patients: patients,
			Consultations: consultations,
		}, nil

	default:
		return nil, errors.New("invalid role")
	}
}

func (s *UserService) ReserveTherapistAgenda(patientUserID, therapistID, agendaID int) error {
	patientID, err := s.Repo.GetPatientIDByUserID(patientUserID)
	if err != nil {
		return err
	}

	agenda, err := s.Repo.GetAgendaByID(agendaID)
	if err != nil {
		return err
	}

	if agenda.Reserved {
		return errors.New("agenda já reservada")
	}

	if agenda.ProfessionalID != therapistID {
		return errors.New("agenda inválida")
	}

	err = s.Repo.MarkAgendaReserved(agendaID, patientID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ReservePsychiatristAgenda(patientUserID, psychiatristID, agendaID int) error {
	patientID, err := s.Repo.GetPatientIDByUserID(patientUserID)
	if err != nil {
		return err
	}

	agenda, err := s.Repo.GetAgendaByID(agendaID)
	if err != nil {
		return err
	}

	if agenda.Reserved {
		return errors.New("agenda já reservada")
	}

	if agenda.ProfessionalID != psychiatristID {
		return errors.New("agenda inválida")
	}

	err = s.Repo.MarkAgendaReserved(agendaID, patientID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ShowConsultation(userID, consultationID int) (models.Consultation, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return models.Consultation{}, err
	}

	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return models.Consultation{}, err
	}

	switch user.Role {

	case "patient":
		patientID, err := s.Repo.GetPatientIDByUserID(userID)
		if err != nil {
			return models.Consultation{}, err
		}
		if c.PatientID != patientID {
			return models.Consultation{}, errors.New("forbidden")
		}

	case "therapist":
		therapistID, err := s.Repo.GetTherapistIDByUserID(userID)
		if err != nil {
			return models.Consultation{}, err
		}
		if c.TherapistID == nil || *c.TherapistID != therapistID {
			return models.Consultation{}, errors.New("forbidden")
		}

	case "psychiatrist":
		psychiatristID, err := s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return models.Consultation{}, err
		}
		if c.PsychiatristID == nil || *c.PsychiatristID != psychiatristID {
			return models.Consultation{}, errors.New("forbidden")
		}

	default:
		return models.Consultation{}, errors.New("invalid role")
	}

	return c, nil
}

func (s *UserService) StartConsultation(userID, agendaID int) (int, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return 0,err
	}

	agenda, err := s.Repo.GetAgendaByID(agendaID)
	if err != nil {
		return 0, err
	}

	if !agenda.Reserved {
		return 0, errors.New("agenda não está reservada")
	}

	if agenda.PatientID == 0 {
		return 0, errors.New("agenda sem paciente")
	}

	err = s.validateAgendaAccess(userID, agenda, user.Role)
	if err != nil {
		return 0, err
	}

	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var consultationID int

	switch user.Role {

	case "therapist":
		price, err := s.Repo.GetTherapistPrice(agenda.ProfessionalID)
		if err != nil {
			return 0, err
		}

		consultationID, err = s.Repo.CreateTherapistConsultation(
			tx,
			agenda.PatientID,
			agenda.ProfessionalID,
			agenda.ID,
			price,
		)
		if err != nil {
			return 0, err
		}

	case "psychiatrist":
		price, err := s.Repo.GetPsychiatristPrice(agenda.ProfessionalID)
		if err != nil {
			return 0, err
		}

		consultationID, err = s.Repo.CreatePsychiatristConsultation(
			tx,
			agenda.PatientID,
			agenda.ProfessionalID,
			agenda.ID,
			price,
		)
		if err != nil {
			return 0, err
		}

	default:
		return 0, errors.New("forbidden")
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return consultationID, nil
}

func (s *UserService) ShowConsultationRoom(userID, consultationID int) (models.ConsultationRoomView, error) {
	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return models.ConsultationRoomView{}, err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return models.ConsultationRoomView{}, err
	}

	if err := s.validateConsultationAccess(userID, c, user.Role); err != nil {
		return models.ConsultationRoomView{}, err
	}

	return models.ConsultationRoomView{
		UserRole: user.Role,
		Consultation: c,
	}, nil
}

func (s *UserService) SaveConsultationRemedy(userID, consultationID int, remedyName, remedyDosage string, remedyQuantity int) error {
	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	switch user.Role {
	case "psychiatrist":
		id, err := s.Repo.GetPsychiatristIDByUserID(userID)
		if err != nil {
			return err
		}

		if c.PsychiatristID == nil || *c.PsychiatristID != id {
			return errors.New("forbidden")
		}

		remedyID, err := s.Repo.InsertRemedy(remedyName, remedyDosage, remedyQuantity)
		if err != nil {
			return err
		}

		return s.Repo.LinkRemedyToConsultation(consultationID, remedyID)

	default:
		return errors.New("forbidden")
	}
}

func (s *UserService) SaveConsultationDiagnosis(userID, consultationID int, diagnosis string) error {
	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Role != "psychiatrist" {
		return errors.New("forbidden")
	}

	id, err := s.Repo.GetPsychiatristIDByUserID(userID)
	if err != nil {
		return err
	}

	if c.PsychiatristID == nil || *c.PsychiatristID != id {
		return errors.New("forbidden")
	}

	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.Repo.UpdateConsultationDiagnosis(tx, consultationID, diagnosis)
	if err != nil {
		return err
	}

	err = s.Repo.UpdatePatientDiagnosis(tx, c.PatientID, diagnosis)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (s *UserService) SaveConsultationBook(userID, consultationID int, author, title string) error {
	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Role != "therapist" {
		return errors.New("forbidden")
	}
	id, err := s.Repo.GetTherapistIDByUserID(userID)
	if err != nil {
		return err
	}

	if c.TherapistID == nil || *c.TherapistID != id {
		return errors.New("forbidden")
	}

	bookID, err := s.Repo.InsertBook(author, title)
	if err != nil {
		return err
	}

	return s.Repo.LinkBookToConsultation(consultationID, bookID)
}

func (s *UserService) SaveConsultationAnnotation(userID, consultationID int, annotation string) error {
	c, err := s.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	err = s.validateConsultationAccess(userID, c, user.Role)
	if err != nil {
		return err
	}

	return s.Repo.UpdateAnnotationConsultation(consultationID, annotation)
}
