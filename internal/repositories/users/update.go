package users

import (
	"database/sql"
	"errors"
)

func (r *UserRepository) CompleteTherapist(userID int, specialty string, description string, price float64) error {
	_, err := r.DB.Exec(`
		UPDATE therapists
		SET specialty = $1, description = $2, price = $3
		WHERE user_id = $4;
	`, specialty, description, price, userID)

	return err
}

func (r *UserRepository) CompletePsychiatrist(userID int, crm string, description string, price float64) error {
	_, err := r.DB.Exec(`
		UPDATE psychiatrists
		SET crm = $1, description = $2, price = $3
		WHERE user_id = $4;
	`, crm, description, price, userID)

	return err
}

func (r *UserRepository) AddTherapistToPatient(patientID int, therapistID int) error {
	_, err := r.DB.Exec(`
		UPDATE patients
		SET therapist_id = $1
		WHERE user_id = $2;
	`, therapistID, patientID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) AddPsychiatristToPatient(patientID int, psychiatristID int) error {
	_, err := r.DB.Exec(`
		UPDATE patients
		SET psychiatrist_id = $1
		WHERE user_id = $2
	`, psychiatristID, patientID)

	return err
}

func (r *UserRepository) UpdateTherapistPrice(userID int, price float64) error {
	_, err := r.DB.Exec(`
		UPDATE therapists 
		SET price = $1 
		WHERE user_id = $2
	`, price, userID)
	return err
}

func (r *UserRepository) UpdatePsychiatristPrice(userID int, price float64) error {
	_, err := r.DB.Exec(`
		UPDATE psychiatrists 
		SET price = $1 
		WHERE user_id = $2
	`, price, userID)
	return err
}

func (r *UserRepository) MarkAgendaReserved(agendaID int, patientID int) error {
	res, err := r.DB.Exec(`
		UPDATE agendas
		SET reserved = true,
		    patient_id = $2
		WHERE id = $1
		AND reserved = false
	`, agendaID, patientID)
	
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("agenda já reservada")
	}

	return nil
}

func (r *UserRepository) UnreserveAgendaPatient(agendaID int, patientID int) error {
	res, err := r.DB.Exec(`
		UPDATE agendas
		SET reserved = false,
		    patient_id = NULL
		WHERE id = $1
		AND patient_id = $2
		AND reserved = true
	`, agendaID, patientID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("agenda não encontrada ou não pertence ao paciente")
	}

	return nil
}

func (r *UserRepository) UnreserveAgendaProfessional(agendaID int, professionalID int) error {
	res, err := r.DB.Exec(`
		UPDATE agendas
		SET reserved = false,
			patient_id = NULL
		WHERE id = $1
		AND professional_id = $2
		AND reserved = true
		
	`, agendaID, professionalID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("agenda não encontrada ou não pertence ao paciente")
	}

	return nil
}

func (r *UserRepository) UpdateAnnotationConsultation(consultationID int, annotation string) error {
	_, err := r.DB.Exec(`
		UPDATE consultations
		SET annotation = $1
		WHERE id = $2
	`, annotation, consultationID)

	return err
}

func (r *UserRepository) UpdateConsultationDiagnosis(tx *sql.Tx, consultationID int, diagnosis string) error {
	_, err := tx.Exec(`
		UPDATE consultations
		SET diagnosis = $1
		WHERE id = $2
	`, diagnosis, consultationID)

	return err
}

func (r *UserRepository) UpdatePatientDiagnosis(tx *sql.Tx, patientID int, diagnosis string) error {
	_, err := tx.Exec(`
		UPDATE patients
		SET current_diagnosis = $1
		WHERE id = $2
	`, diagnosis, patientID)

	return err
}