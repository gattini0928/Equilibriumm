package validators

import (
	"errors"
)

func ValidateAge(age int, role string) error {
	if age < 21 {
		return errors.New("Você precisa ser maior de 21 anos para exercer a função de terapeuta ou psiquiatra")
	}

	if age > 90 {
		return errors.New("Idade inválida")
	}
	
	return nil
}