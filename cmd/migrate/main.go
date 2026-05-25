package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"

	"github.com/gattini0928/Equilibrium/internal/db"
)

func main() {
	godotenv.Load()
	conn := db.Connect()
	defer conn.Close()

	drop := `
	DROP TABLE IF EXISTS consultation_remedies CASCADE;
	DROP TABLE IF EXISTS consultation_books CASCADE;
	DROP TABLE IF EXISTS consultations CASCADE;
	DROP TABLE IF EXISTS agendas CASCADE;
	DROP TABLE IF EXISTS remedies CASCADE;
	DROP TABLE IF EXISTS books CASCADE;
	DROP TABLE IF EXISTS patients CASCADE;
	DROP TABLE IF EXISTS therapists CASCADE;
	DROP TABLE IF EXISTS psychiatrists CASCADE;
	DROP TABLE IF EXISTS users CASCADE;
	`

	_, err := conn.Exec(drop)
	if err != nil {
		log.Fatal(err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			email VARCHAR(200) UNIQUE,
			password VARCHAR(255),
			age INTEGER,
			cpf VARCHAR(11) UNIQUE,
			role VARCHAR(20),
			image TEXT
		);

		CREATE TABLE IF NOT EXISTS therapists (
			id SERIAL PRIMARY KEY,
			user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			specialty TEXT,
			description TEXT,
			price REAL
		);

		CREATE TABLE IF NOT EXISTS psychiatrists (
			id SERIAL PRIMARY KEY,
			user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			crm TEXT,
			description TEXT,
			price REAL
		);

		CREATE TABLE IF NOT EXISTS patients (
			id SERIAL PRIMARY KEY,
			user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			therapist_id INTEGER REFERENCES therapists(id),
			psychiatrist_id INTEGER REFERENCES psychiatrists(id),
			current_diagnosis TEXT
		);

		CREATE TABLE IF NOT EXISTS agendas (
			id SERIAL PRIMARY KEY,
			professional_id INTEGER,
			professional_role TEXT NOT NULL CHECK (professional_role IN ('therapist', 'psychiatrist')),
			patient_id INTEGER,
			date TIMESTAMP NOT NULL, 
			reserved BOOLEAN
		);

		CREATE TABLE IF NOT EXISTS consultations (
			id SERIAL PRIMARY KEY,
			patient_id INTEGER REFERENCES patients(id) ON DELETE CASCADE,
			therapist_id INTEGER REFERENCES therapists(id) ON DELETE CASCADE,
			psychiatrist_id INTEGER REFERENCES psychiatrists(id) ON DELETE CASCADE,
			CHECK (
				(therapist_id IS NOT NULL AND psychiatrist_id IS NULL)
				OR
				(therapist_id IS NULL AND psychiatrist_id IS NOT NULL)
			),
			date TIMESTAMP NOT NULL,
			price REAL,
			annotation TEXT NOT NULL DEFAULT '',
			diagnosis TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			author TEXT,
			title TEXT
		);

		CREATE TABLE IF NOT EXISTS remedies (
			id SERIAL PRIMARY KEY,
			name TEXT,
			dosage TEXT,
			quantity INTEGER
		);

		CREATE TABLE IF NOT EXISTS consultation_books (
			consultation_id INTEGER REFERENCES consultations(id) ON DELETE CASCADE,
			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
			UNIQUE (consultation_id, book_id)
		);

		CREATE TABLE IF NOT EXISTS consultation_remedies (
			consultation_id INTEGER REFERENCES consultations(id) ON DELETE CASCADE,
			remedy_id INTEGER REFERENCES remedies(id) ON DELETE CASCADE,
			UNIQUE (consultation_id, remedy_id)
		);

	`
	
	_, err = conn.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tabelas criadas com sucesso 🚀")
}