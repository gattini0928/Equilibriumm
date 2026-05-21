package users

import (
	"errors"
)

func (r *UserRepository) DeleteAgenda(userID int, agendaID int) error {
	query := `
		DELETE FROM agendas a
		WHERE a.id = $1
		AND a.professional_id = $2
		AND a.reserved = false
		AND NOT EXISTS(
			SELECT 1 FROM consultations c
			WHERE c.agenda_id = a.id
		)
	`
	res, err := r.DB.Exec(query, agendaID, userID)
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





