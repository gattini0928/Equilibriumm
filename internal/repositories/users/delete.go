package users

import (
	"database/sql"
	"errors"
)

func (r *UserRepository) DeleteAgendaTX(tx *sql.Tx, agendaID int, professionalID int, professionalRole string, patientID int) error {
	query := `
		DELETE FROM agendas
		WHERE id = $1
		AND professional_id = $2
		AND professional_role = $3
		AND patient_id = $4
		AND reserved = true
	`
	res, err := tx.Exec(query, agendaID, professionalID, professionalRole, patientID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return errors.New("agenda não encontrada ou não pertence ao usuário")
	}

	return nil
}

func (r *UserRepository) DeleteAgenda(agendaID int, professionalID int, professionalRole string) error {
	query := `
		DELETE FROM agendas
		WHERE id = $1
		AND professional_id = $2
		AND professional_role = $3
		AND reserved = false
	`
	res, err := r.DB.Exec(query, agendaID, professionalID, professionalRole)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return errors.New("agenda não encontrada ou não pertence ao usuário")
	}

	return nil
}






