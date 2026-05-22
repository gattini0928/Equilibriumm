package validators

import ("time" 
	"errors"
)

func ValidateDate(date time.Time) error {
	if date.Before(time.Now()) {
		return errors.New("data inválida")
	}

	return nil
}

func ValidateHour(hour string) error {
	t, err := time.Parse("15:04", hour)
	if err != nil {
		return errors.New("horário inválido, use 07:00")
	}

	if t.Hour() < 7 || t.Hour() > 21 {
		return errors.New("horário deve ser entre 07:00 e 21:00")
	}

	return nil
}