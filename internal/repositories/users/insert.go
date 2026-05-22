package users

import (
	"database/sql"
	"errors"

	"github.com/gattini0928/Equilibrium/internal/models"
)

func (r *UserRepository) CreateUserWithProfile(
	user *models.User,
	patient *models.Patient,
	therapist *models.Therapist,
	psychiatrist *models.Psychiatrist,
) error {

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	var userID int

	err = tx.QueryRow(`
		INSERT INTO users (name, email, password, age, cpf, role, image)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`, user.Name, user.Email, user.Password, user.Age, user.Cpf, user.Role, user.Image).Scan(&userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	user.ID = userID

	switch user.Role {

	case "patient":
		_, err = tx.Exec(`
			INSERT INTO patients (user_id)
			VALUES ($1);
		`, userID)

	case "therapist":
		_, err = tx.Exec(`
			INSERT INTO therapists (user_id)
			VALUES ($1);
		`, userID)

	case "psychiatrist":
		_, err = tx.Exec(`
			INSERT INTO psychiatrists (user_id)
			VALUES ($1);
		`, userID)

	default:
		tx.Rollback()
		return errors.New("invalid role")
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) InsertAgenda(professionalID int, professionalRole string,day int, month int, hour string) (models.Agenda, error) {
	var agenda models.Agenda

	query := `
		INSERT INTO agendas (professional_id, professional_role, day, month, hour, reserved)
		VALUES ($1, $2, $3, $4, $5, false)
		RETURNING id, professional_id, professional_role, day, month, hour, reserved;
	`

	err := r.DB.QueryRow(query, professionalID, professionalRole, day, month, hour).Scan(
		&agenda.ID,
		&agenda.ProfessionalID,
		&agenda.ProfessionalRole,
		&agenda.Day,
		&agenda.Month,
		&agenda.Hour,
		&agenda.Reserved,
	)

	if err != nil {
		return models.Agenda{}, err
	}

	return agenda, nil
}

func (r *UserRepository) CreateTherapistConsultation(tx *sql.Tx,  patientID, therapistID, agendaID int, price float64) (int, error) {
	var consultationID int
	query := `
		INSERT INTO consultations (patient_id, therapist_id, agenda_id, date, price)
		VALUES ($1, $2, $3, NOW(), $4) RETURNING id
		`
	err := tx.QueryRow(query, patientID, therapistID, agendaID, price).Scan(&consultationID)
	return consultationID, err
}

func (r *UserRepository) CreatePsychiatristConsultation(tx *sql.Tx, patientID, psychiatristID, agendaID int, price float64) (int, error) {
	var consultationID int
	query := `
		INSERT INTO consultations (patient_id, psychiatrist_id, agenda_id, date, price)
		VALUES ($1, $2, $3, NOW(), $4) RETURNING id
		`
	err := tx.QueryRow(query, patientID, psychiatristID, agendaID, price).Scan(&consultationID)
	return consultationID, err
}

func (r *UserRepository) InsertRemedy(name, dosage string, quantity int) (int, error) {
	var id int

	err := r.DB.QueryRow(`
		INSERT INTO remedies (name, dosage, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, dosage, quantity).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *UserRepository) LinkRemedyToConsultation(consultationID, remedyID int) error {
	_, err := r.DB.Exec(`
		INSERT INTO consultation_remedies (consultation_id, remedy_id)
		VALUES ($1, $2)
	`, consultationID, remedyID)

	return err
}

func(r *UserRepository) InsertBook(author, title string) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO books (author, title)
		VALUES ($1, $2) RETURNING id
	`, author, title).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (r *UserRepository) LinkBookToConsultation(consultationID, bookID int) error {
	_, err := r.DB.Exec(`
		INSERT INTO consultation_books (consultation_id, book_id)
		VALUES ($1, $2)
	`, consultationID, bookID)

	return err
}

