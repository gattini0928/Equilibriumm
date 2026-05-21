package models

import "time"

type User struct {
	ID int
	Name string
	Email string
	Password string
	Age int
	Cpf string 
	Role string
	Image string
}

type UserPerfil struct {
	ID int
	Name string
	Email string
	Age int
	Role string
	Image string
	CurrentDiagnosis string
	Specialty   string
	CRM string
	Description string
	Price float64
}

type Patient struct {
	ID               int 
	UserID           int   
	TherapistID      *int  
	PsychiatristID   *int 
	CurrentDiagnosis string
	Therapist *TherapistInfo
	Psychiatrist *PsychiatristInfo 
}

type Therapist struct {
	ID          int    
	UserID      int    
	Specialty   string
	Description string
	Price float64
}

type Psychiatrist struct {
	ID     int   
	UserID int   
	CRM    string
	Description string
	Price float64
}

type TherapistInfo struct {
	ID int
	Name string
	Email string
	Age int
	Image string 
	Specialty string
	Description string
}

type PsychiatristInfo struct {
	ID int 
	Name string
	Email string
	Age int
	Image string
	Description string
}

type Consultation struct {
	ID               int
	PatientID        int
	PatientName string
	DoctorName string
	DoctorRole string
	TherapistID   *int    
	PsychiatristID *int  
	Date             time.Time 
	Price            float64 
	Annotation       string
	Diagnosis string  
	AgendaID int 
	Books    []Book
	Remedies []Remedy
}

type ConsultationRoomView struct {
	UserRole     string
	Consultation Consultation
}

type Book struct {
	ID     int  
	Author string
	Title  string
}

type Remedy struct {
	ID       int   
	Name     string
	Dosage   string
	Quantity int   
}

type Agenda struct {
	ID             int
	ProfessionalID int
	ProfessionalRole string
	PatientID 		int
	PatientName    string
	Day            int  
	Month          int 
	Hour           string
	Reserved       bool  
}

type DoctorWithUser struct {
	ID          int
	Name        string
	Email       string
	Image       string
	Age         int
	Specialty   string
	CRM 		string
	Description string
	Price float64
}

type PatientWithUser struct {
	ID          int
	Name        string
	Email       string
	Image       string
	Age         int
	CurrentDiagnosis string
	Books       []Book
	Remedies []Remedy
}

type PatientDashboard struct {
	Perfil UserPerfil
	Role string
	Therapist *DoctorWithUser
	Psychiatrist *DoctorWithUser
	AgendasReserved []Agenda
	Consultations []Consultation
}

type DoctorDashboard struct {
	Perfil   UserPerfil
	Role string 
	Agendas  []Agenda
	AgendasReserved []Agenda     
	Patients []PatientWithUser
	Consultations []Consultation
}

