package product

import (
	"errors"
	"regexp"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyRequestProduct struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (p *Product) ValidateBusinessRule() error {
	if p.Id == 0 {
		return errors.New("invalid field Id")
	}
	if p.Name == "" {
		return errors.New("invalid field Name")
	}
	if p.Quantity == 0 {
		return errors.New("invalid field Quantity")
	}
	if p.CodeValue == "" {
		return errors.New("invalid field CodeValue")
	}

	if p.Expiration == "" {
		return errors.New("invalid field Expiration")
	}

	expirationFormat, _ := regexp.MatchString(`^(0[1-9]|[1-2][0-9]|3[0-1])/(0[1-9]|1[0-2])/\d{4}$`, p.Expiration)
	if !expirationFormat {
		return errors.New("invalid format Expiration")
	}

	if p.Price == 0.0 {
		return errors.New("invalid field Price")
	}
	return nil
}
