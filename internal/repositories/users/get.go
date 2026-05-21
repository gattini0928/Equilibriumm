package users

import (
	"database/sql"
	"errors"

	"github.com/gattini0928/Equilibrium/internal/models"
)

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, email, password, age, cpf, role, image 
		FROM users 
		WHERE email = $1
	`, email)
	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Age,
		&user.Cpf,
		&user.Role,
		&user.Image,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User

	row := r.DB.QueryRow(`
		SELECT id, name, email, age, image, role
		FROM users 
		WHERE id = $1`, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Image,
		&user.Role,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetPatientIDByUserID(userID int) (int, error) {
	var patientID int

	err := r.DB.QueryRow(`
		SELECT id
		FROM patients
		WHERE user_id = $1	
	`, userID).Scan(&patientID)

	if err != nil {
		return 0, err
	}

	return patientID, nil
}

func (r *UserRepository) GetPatientPerfil(userID int) (models.UserPerfil, error) {
	query := `
		SELECT u.id, u.name, u.email, u.age, u.image, p.current_diagnosis
		FROM patients p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
	`
	var patient models.UserPerfil
	var currentDiagnosis sql.NullString

	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&patient.ID,
		&patient.Name,
		&patient.Email,
		&patient.Age,
		&patient.Image,
		&currentDiagnosis,
	)

	if err != nil {
		return models.UserPerfil{}, err
	}

	if currentDiagnosis.Valid {
		patient.CurrentDiagnosis = currentDiagnosis.String
	}

	return patient, nil
}

func(r *UserRepository) GetPatientReservedAgendas(patientID int) ([]models.Agenda, error) {
	query := `
		SELECT
			id,
			professional_id,
			day,
			month,
			hour,
			reserved
		FROM agendas
		WHERE patient_id = $1
		AND reserved = true
	`

	rows, err := r.DB.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.ProfessionalID,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetTherapistPerfil(userID int) (models.UserPerfil, error) {
	query := `
		SELECT u.id, u.name, u.email, u.age, u.image, u.role, t.specialty, t.description, t.price
		FROM therapists t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1
	`
	var therapist models.UserPerfil
	var price sql.NullFloat64

	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&therapist.ID,
		&therapist.Name,
		&therapist.Email,
		&therapist.Age,
		&therapist.Image,
		&therapist.Role,
		&therapist.Specialty,
		&therapist.Description,
		&price,
	)

	if err != nil {
		return models.UserPerfil{}, err
	}

	if price.Valid {
		therapist.Price = price.Float64
	}

	return therapist, nil
}

func (r *UserRepository) GetPsychiatristPerfil(userID int) (models.UserPerfil, error) {
	query := `
		SELECT u.id, u.name, u.email, u.age, u.image, u.role, p.crm, p.description, p.price
		FROM psychiatrists p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
	`
	var psychiatrist models.UserPerfil
	var price sql.NullFloat64

	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&psychiatrist.ID,
		&psychiatrist.Name,
		&psychiatrist.Email,
		&psychiatrist.Age,
		&psychiatrist.Image,
		&psychiatrist.Role,
		&psychiatrist.CRM,
		&psychiatrist.Description,
		&price,
	)

	if err != nil {
		return models.UserPerfil{}, err
	}

	if price.Valid {
		psychiatrist.Price = price.Float64
	}

	return psychiatrist, nil
}

func (r *UserRepository) GetAllTherapists() ([]models.DoctorWithUser, error) {
	query := `
		SELECT 
			t.id,
			u.name,
			u.email,
			u.age,
			u.image,
			t.specialty,
			t.description
		FROM therapists t
		JOIN users u ON t.user_id = u.id
		WHERE t.description IS NOT NULL AND t.specialty IS NOT NULL 
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var therapists []models.DoctorWithUser

	for rows.Next() {
		var t models.DoctorWithUser
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Age,
			&t.Image,
			&t.Specialty,
			&t.Description,
		)

		if err != nil {
			return nil, err
		}

		therapists = append(therapists, t)
	}

	if therapists == nil {
    	therapists = []models.DoctorWithUser{}
	}

	return therapists, nil
}

func (r *UserRepository) GetAllPsychiatrists() ([]models.DoctorWithUser, error) {
	query := `
		SELECT 
			p.id,
			u.name,
			u.email,
			u.age,
			u.image,
			p.description,
			p.crm
		FROM psychiatrists p
		JOIN users u ON p.user_id = u.id
		WHERE p.description IS NOT NULL AND p.crm IS NOT NULL
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var psychiatrists []models.DoctorWithUser

	for rows.Next() {
		var p models.DoctorWithUser
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Email,
			&p.Age,
			&p.Image,
			&p.Description,
			&p.CRM,
		)
		if err != nil {
			return nil, err
		}

		psychiatrists = append(psychiatrists, p)
	}

	if psychiatrists == nil {
    	psychiatrists = []models.DoctorWithUser{}
	}

	return psychiatrists, nil
}

func (r *UserRepository) GetTherapistById(userID int) (models.DoctorWithUser, error) {
	var therapist models.DoctorWithUser
	query := 
		`
		SELECT 
			t.id,
			u.name,
			u.email,
			u.age,
			u.image,
			t.specialty,
			t.description,
			t.price
		FROM therapists t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1;	
		`
	var price sql.NullFloat64

	row := r.DB.QueryRow(query, userID)
		err := row.Scan(
		&therapist.ID,
		&therapist.Name,
		&therapist.Email,
		&therapist.Age,
		&therapist.Image,
		&therapist.Specialty,
		&therapist.Description,
		&price,
	)

	if err != nil {
		return models.DoctorWithUser{}, err
	}

	if price.Valid {
		therapist.Price = price.Float64
	}

	return therapist, nil
}

func (r *UserRepository) GetTherapistAgenda(therapistID int) ([]models.Agenda, error ) {
	query := `
		SELECT id, day, month, hour, reserved
		FROM agendas
		WHERE professional_id = $1
		AND reserved = false
	`
	rows, err := r.DB.Query(query,therapistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetTherapistPrivateAgenda(userID int) ([]models.Agenda, error ) {
	query := `
		SELECT id, day, month, hour, reserved
		FROM agendas
		WHERE professional_id = $1
	`
	rows, err := r.DB.Query(query,userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetTherapistReservedAgendas(therapistID int) ([]models.Agenda ,error) {
	query := `
		SELECT
			a.id,
			a.patient_id,
			u.name,
			a.day,
			a.month,
			a.hour,
			a.reserved
		FROM agendas a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = p.user_id
		WHERE a.professional_id = $1
		AND a.reserved = true
	`

	rows, err := r.DB.Query(query, therapistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.PatientID,
			&a.PatientName,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetPsychiatristById(userID int) (models.DoctorWithUser, error) {
	var psychiatrist models.DoctorWithUser

	query := `
		SELECT p.id,
			u.name,
			u.email,
			u.age,
			u.image,
			p.crm,
			p.description,
			p.price
		FROM psychiatrists p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1;
	`

	var price sql.NullFloat64
	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&psychiatrist.ID,
		&psychiatrist.Name,
		&psychiatrist.Email,
		&psychiatrist.Age,
		&psychiatrist.Image,
		&psychiatrist.CRM,
		&psychiatrist.Description,
		&price,
	)

	if err != nil {
		return models.DoctorWithUser{}, err
	}

	if price.Valid {
		psychiatrist.Price = price.Float64
	}

	return psychiatrist, nil
}

func (r *UserRepository) GetPsychiatristAgenda(psychiatristID int) ([]models.Agenda, error ) {
	query := `
		SELECT id, day, month, hour, reserved
		FROM agendas
		WHERE professional_id = $1
		AND reserved = false
	`
	rows, err := r.DB.Query(query,psychiatristID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetPsychiatristPrivateAgenda(userID int) ([]models.Agenda, error ) {
	query := `
		SELECT
			a.id,
			a.patient_id,
			u.name,
			a.day,
			a.month,
			a.hour,
			a.reserved
		FROM agendas a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = p.user_id
		WHERE a.professional_id = $1
		AND a.reserved = true
	`
	rows, err := r.DB.Query(query,userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.PatientID,
			&a.PatientName,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetPsychiatristReservedAgendas(psychiatristID int) ([]models.Agenda ,error) {
	query := `
		SELECT id, day, month, hour, reserved
		FROM agendas
		WHERE professional_id = $1
		AND reserved = true
	`

	rows, err := r.DB.Query(query, psychiatristID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var a models.Agenda
		err := rows.Scan(
			&a.ID,
			&a.Day,
			&a.Month,
			&a.Hour,
			&a.Reserved,
		)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, a)
	}

	if agendas == nil {
		return []models.Agenda{}, nil
	}

	return agendas, nil
}

func (r *UserRepository) GetTherapistPatient(patiendID int) (models.PatientWithUser, error) {
	var patient models.PatientWithUser
	query := `
	SELECT 
		p.id,
		u.name, u.email, u.age, u.image,
		p.current_diagnosis,
		b.id, b.author, b.title
		FROM patients p
		JOIN users u ON u.id = p.user_id
		JOIN consultations c ON c.patient_id = p.id
		JOIN consultation_books cb ON cb.consultation_id = c.id
		JOIN books b ON b.id = cb.book_id
		WHERE p.user_id = $1;
		`

	rows, err := r.DB.Query(query, patiendID)
	if err != nil {
		return models.PatientWithUser{}, err
	}

	defer rows.Close()

	var books []models.Book
	for rows.Next(){
		var b models.Book
		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.Email,
			&patient.Age,
			&patient.Image,
			&patient.CurrentDiagnosis,
			&b.ID,
			&b.Author,
			&b.Title,
		)
		if err != nil {
			return models.PatientWithUser{}, err
		}
		books = append(books, b)
	}

	patient.Books = books
	return patient, nil
}

func (r *UserRepository) GetPsychiatristPatient(patiendID int) (models.PatientWithUser, error) {
	var patient models.PatientWithUser
	query := `
		SELECT p.id, u.name, u.email, u.age, u.image, p.current_diagnosis,
		r.id, r.name, r.dosage, r.quantity
		FROM patients p
		JOIN users u ON u.id = p.id
		JOIN consultation c ON c.patient_id = p.id
		JOIN consultation_remedies cr ON consultation.id = cr.remedy_id
		JOIN remedies r ON r.id = cr.id
		WHERE user_id = $1;
	`

	rows, err := r.DB.Query(query,patiendID)
	if err != nil {
		return models.PatientWithUser{}, err
	}
	defer rows.Close()

	var remedies []models.Remedy

	for rows.Next() {
		var r models.Remedy
		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.Email,
			&patient.Age,
			&patient.Image,
			&patient.CurrentDiagnosis,
			&r.ID,
			&r.Name,
			&r.Dosage,
			&r.Quantity,
		)
		if err != nil {
			return models.PatientWithUser{}, err
		}

		remedies = append(remedies, r)
	}

	patient.Remedies = remedies

	return patient, nil
}

 func (r *UserRepository) GetTherapistPatients(doctorID int) ([]models.PatientWithUser, error) {
	query := `
	SELECT p.id, u.name, u.email, u.age, u.image, p.current_diagnosis
	FROM patients p
	JOIN users u ON u.id = p.user_id
	JOIN therapists t ON t.id = p.therapist_id
	WHERE t.user_id = $1
	`

	rows, err := r.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.PatientWithUser
	var currentDiagnosis sql.NullString

	for rows.Next() {
		var p models.PatientWithUser
		err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Age, &p.Image, &currentDiagnosis)
		if err != nil {
			return nil, err
		}

		if currentDiagnosis.Valid {
			p.CurrentDiagnosis = currentDiagnosis.String
		}

		patients = append(patients, p)
	}

	if patients == nil {
		return []models.PatientWithUser{}, nil
	}

	return patients, nil
}

 func (r *UserRepository) GetPsychiatristPatients(doctorID int) ([]models.PatientWithUser, error) {
	query := `
	SELECT p.id, u.name, u.email, u.age, u.image, p.current_diagnosis
	FROM patients p
	JOIN users u ON u.id = p.user_id
	JOIN psychiatrists ps ON p.id = p.psychiatrist_id
	WHERE ps.user_id = $1
	`

	rows, err := r.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.PatientWithUser

	for rows.Next() {
		var p models.PatientWithUser
		err := rows.Scan(
			&p.ID, &p.Name, &p.Email, &p.Age, &p.Image, 
			&p.CurrentDiagnosis)

		if err != nil {
			return nil, err
		}

		patients = append(patients, p)
	}

	
	if patients == nil {
		return []models.PatientWithUser{}, nil
	}

	return patients, nil
}

func (r *UserRepository) GetPatientTherapist(userID int) (*models.DoctorWithUser, error) {
	var therapist models.DoctorWithUser

	query := `
	SELECT 
		t.id,
		u.name,
		u.email,
		u.age,
		u.image,
		t.specialty,
		t.description
	FROM patients p
	JOIN therapists t ON t.id = p.therapist_id
	JOIN users u ON u.id = t.user_id
	WHERE p.user_id = $1;
	`

	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&therapist.ID,
		&therapist.Name,
		&therapist.Email,
		&therapist.Age,
		&therapist.Image,
		&therapist.Specialty,
		&therapist.Description,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &therapist, nil
}

func (r *UserRepository) GetPatientPsychiatrist(userID int) (*models.DoctorWithUser, error) {
	var psychiatrist models.DoctorWithUser

	query := `
	SELECT 
		p.id,
		u.name,
		u.email,
		u.age,
		u.image,
		p.description
	FROM patients pat
	JOIN psychiatrists p ON p.id = pat.psychiatrist_id
	JOIN users u ON u.id = p.user_id
	WHERE pat.user_id = $1;
	`

	row := r.DB.QueryRow(query, userID)
	err := row.Scan(
		&psychiatrist.ID,
		&psychiatrist.Name,
		&psychiatrist.Email,
		&psychiatrist.Age,
		&psychiatrist.Image,
		&psychiatrist.Description,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &psychiatrist, nil
}

func (r *UserRepository) GetAgendaByID(agendaID int) (models.Agenda, error) {
	var a models.Agenda
	var patientID sql.NullInt64

	err := r.DB.QueryRow(`
		SELECT id, professional_id, patient_id, day, month, hour, reserved
		FROM agendas
		WHERE id = $1
	`, agendaID).Scan(
		&a.ID,
		&a.ProfessionalID,
		&patientID,
		&a.Day,
		&a.Month,
		&a.Hour,
		&a.Reserved,
	)

	if patientID.Valid {
		a.PatientID = int(patientID.Int64)
	}

	if err != nil {
		return models.Agenda{}, err
	}

	return a, nil
}

func (r *UserRepository) GetTherapistPrice(therapistID int) (float64, error) {
	var price sql.NullFloat64

	err := r.DB.QueryRow(`
		SELECT price
		FROM therapists
		WHERE id = $1
	`, therapistID).Scan(&price)

	if err != nil {
		return 0, err
	}

	if !price.Valid {
		return 0, errors.New("therapist sem preço definido")
	}

	return price.Float64, nil
}

func (r *UserRepository) GetPsychiatristPrice(psychiastridID int) (float64, error) {
	var price sql.NullFloat64

	err := r.DB.QueryRow(
		`SELECT price
			FROM psychiatrists WHERE id = $1
			`, psychiastridID).Scan(&price)

	if err != nil {
		return 0, err
	}

	if !price.Valid {
		return 0, errors.New("therapist sem preço definido")
	}

	return price.Float64, nil
}

func (r *UserRepository) GetTherapistIDByUserID(userID int) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		SELECT id 
		FROM therapists
		WHERE user_id = $1
	`, userID).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetPsychiatristIDByUserID(userID int) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		SELECT id 
		FROM psychiatrists
		WHERE user_id = $1
	`, userID).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (r *UserRepository) GetPatientConsultations(patientID int) ([]models.Consultation,error) {
// 	query := `
// 		SELECT id, patient_id, therapist_id, psychiatrist_id, date, price, annotation, agenda_id
// 		FROM consultations
// 		WHERE patient_id = $1
// 	`
// 	rows, err := r.DB.Query(query, patientID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var consultations []models.Consultation
// 	for rows.Next() {
// 		var consultation models.Consultation
// 		err := rows.Scan(
// 			&consultation.ID,
// 			&consultation.PatientID,
// 			&consultation.TherapistID,
// 			&consultation.PsychiatristID,
// 			&consultation.Date,
// 			&consultation.Price,
// 			&consultation.Annotation,
// 			&consultation.AgendaID,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		consultations = append(consultations, consultation)
// 	}

// 	if consultations == nil {
// 		return []models.Consultation{}, nil
// 	}

// 	return consultations, nil

// }

func (r *UserRepository) GetPatientConsultations(patientID int) ([]models.Consultation, error) {
	query := `
	SELECT 
		c.id,
		c.patient_id,

		CASE
			WHEN c.therapist_id IS NOT NULL THEN tu.name
			WHEN c.psychiatrist_id IS NOT NULL THEN pu.name
		END as doctor_name,

		CASE
			WHEN c.therapist_id IS NOT NULL THEN 'therapist'
			WHEN c.psychiatrist_id IS NOT NULL THEN 'psychiatrist'
		END as doctor_role,

			c.date,
			c.price,
			c.annotation,
			c.diagnosis,
			c.agenda_id

		FROM consultations c

		LEFT JOIN therapists t ON t.id = c.therapist_id
		LEFT JOIN users tu ON tu.id = t.user_id

		LEFT JOIN psychiatrists p ON p.id = c.psychiatrist_id
		LEFT JOIN users pu ON pu.id = p.user_id

		WHERE c.patient_id = $1`

 	rows, err := r.DB.Query(query, patientID)
 	if err != nil {
 		return nil, err
 	}
 	defer rows.Close()

 	var consultations []models.Consultation
 	for rows.Next() {
 		var consultation models.Consultation
 		err := rows.Scan(
			&consultation.ID,
			&consultation.PatientID,
			&consultation.DoctorName,
			&consultation.DoctorRole,
			&consultation.Date,
			&consultation.Price,
			&consultation.Annotation,
			&consultation.Diagnosis,
			&consultation.AgendaID,
		)
		if err != nil {
			return nil, err
	}

	consultations = append(consultations, consultation)
	}

	if consultations == nil {
		return []models.Consultation{}, nil
}

	return consultations, nil
}

func (r *UserRepository) GetTherapistConsultations(therapistID int) ([]models.Consultation,error) {
	query := `
		SELECT 
			c.id,
			c.patient_id,
			u.name,
			c.date
		FROM consultations c
		JOIN patients p ON p.id = c.patient_id
		JOIN users u ON u.id = p.user_id
		WHERE c.therapist_id = $1
	`
	rows, err := r.DB.Query(query, therapistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultations []models.Consultation
	for rows.Next() {
		var consultation models.Consultation
		err := rows.Scan(
			&consultation.ID,
			&consultation.PatientID,
			&consultation.PatientName,
			&consultation.Date,
		)
		if err != nil {
			return nil, err
		}

		consultations = append(consultations, consultation)
	}

	if consultations == nil {
		return []models.Consultation{}, nil
	}

	return consultations, nil

}

func (r *UserRepository) GetPsychiatristConsultations(psychiatristID int) ([]models.Consultation,error) {
	query := `
		SELECT
			c.id,
			c.patient_id,
			u.name,
			c.date
		FROM consultations c
		JOIN patients p ON p.id = c.patient_id
		JOIN users u ON u.id = p.user_id
		WHERE c.psychiatrist_id = $1
	`
	rows, err := r.DB.Query(query, psychiatristID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultations []models.Consultation
	for rows.Next() {
		var consultation models.Consultation
		err := rows.Scan(
			&consultation.ID,
			&consultation.PatientID,
			&consultation.PatientName,
			&consultation.Date,
		)
		if err != nil {
			return nil, err
		}

		consultations = append(consultations, consultation)
	}

	if consultations == nil {
		return []models.Consultation{}, nil
	}

	return consultations, nil

}

func (r *UserRepository) GetConsultationByID(consultationID int) (models.Consultation, error) {
	var c models.Consultation

	err := r.DB.QueryRow(`
		SELECT 
			c.id,
			c.patient_id,
			pu.name,

			c.therapist_id,
			c.psychiatrist_id,

			COALESCE(tu.name, psu.name) AS doctor_name,

			CASE
				WHEN c.therapist_id IS NOT NULL THEN 'therapist'
				WHEN c.psychiatrist_id IS NOT NULL THEN 'psychiatrist'
			END AS doctor_role,

			c.date,
			c.price,
			c.annotation,
			c.diagnosis
		FROM consultations c
		JOIN patients p ON p.id = c.patient_id
		JOIN users pu ON pu.id = p.user_id

		LEFT JOIN therapists t ON t.id = c.therapist_id
		LEFT JOIN users tu ON tu.id = t.user_id

		LEFT JOIN psychiatrists ps ON ps.id = c.psychiatrist_id
		LEFT JOIN users psu ON psu.id = ps.user_id

		WHERE c.id = $1
	`, consultationID).Scan(
		&c.ID,
		&c.PatientID,
		&c.PatientName,

		&c.TherapistID,
		&c.PsychiatristID,

		&c.DoctorName,
		&c.DoctorRole,

		&c.Date,
		&c.Price,
		&c.Annotation,
		&c.Diagnosis,
	)
	
	if err != nil {
		return models.Consultation{}, err
	}

	bookRows, err := r.DB.Query(`
		SELECT b.id, b.author, b.title
		FROM consultation_books cb
		JOIN books b ON b.id = cb.book_id
		WHERE cb.consultation_id = $1
	`, consultationID)
	if err != nil {
		return models.Consultation{}, err
	}
	defer bookRows.Close()

	for bookRows.Next() {
		var b models.Book
		err := bookRows.Scan(&b.ID, &b.Author, &b.Title)
		if err != nil {
			return models.Consultation{}, err
		}
		c.Books = append(c.Books, b)
	}

	remedyRows, err := r.DB.Query(`
		SELECT r.id, r.name, r.dosage, r.quantity
		FROM consultation_remedies cr
		JOIN remedies r ON r.id = cr.remedy_id
		WHERE cr.consultation_id = $1
	`, consultationID)
	if err != nil {
		return models.Consultation{}, err
	}
	defer remedyRows.Close()

	for remedyRows.Next() {
		var rm models.Remedy
		err := remedyRows.Scan(&rm.ID, &rm.Name, &rm.Dosage, &rm.Quantity)
		if err != nil {
			return models.Consultation{}, err
		}
		c.Remedies = append(c.Remedies, rm)
	}

	return c, nil
}

