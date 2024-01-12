package service

import (
	"ejercicio1/internal"
	"errors"
	"fmt"
	"regexp"
)

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

type ProductDefault struct {
	rp internal.ProductRepository
	//external services ...
}

func (p *ProductDefault) GetAll() (products []internal.Product, err error) {
	//external service
	// ...
	//buisness logic
	// ....

	products, err = p.rp.GetAll()
	return
}

func (p *ProductDefault) GetById(id int) (product internal.Product, err error) {
	//external service
	// ...

	//buisness logic
	if id == 0 {
		err = fmt.Errorf("%w: Id", internal.ErrFieldRequire)
		return
	}

	product, err = p.rp.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			err = fmt.Errorf("%w", internal.ErrProductNotFound)
			return
		}
	}
	return
}

func (p *ProductDefault) GetPriceGreaterThan(priceGt float64) (products []internal.Product, err error) {
	//external service
	// ....

	//buisness logic
	if priceGt == 0 {
		err = fmt.Errorf("%w: price", internal.ErrInvalidFormat)
		return
	}

	products, err = p.rp.GetPriceGreaterThan(priceGt)
	if err != nil {
		if errors.Is(err, internal.ErrProductNotFound) {
			err = fmt.Errorf("%w", internal.ErrProductNotFound)
			return
		}
		err = fmt.Errorf("unknown problem")
		return
	}

	return
}

func (p *ProductDefault) Save(product *internal.Product) error {
	//external service

	//buisness logic
	if err := validateBuisnessLogic(*product); err != nil {
		return err
	}
	//save product
	err := p.rp.Save(product)
	if err != nil {
		switch err {
		case internal.ErrProductAlredyExist:
			err = fmt.Errorf("%w: title", internal.ErrProductAlredyExist)
		}
		return err
	}
	return nil
}

func (p *ProductDefault) Update(product *internal.Product) error {
	//buisness logic
	if err := validateBuisnessLogic(*product); err != nil {
		return err
	}
	err := p.rp.Update(product)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			return fmt.Errorf("%w: id", internal.ErrProductNotFound)
		default:
			return fmt.Errorf("unkown error")
		}
	}
	return nil
}

func (p *ProductDefault) Delete(id int) error {
	if id == 0 {
		return fmt.Errorf("%w: Id", internal.ErrFieldRequire)

	}

	err := p.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			return fmt.Errorf("%w: id", internal.ErrProductNotFound)
		default:
			return fmt.Errorf("unkown error")
		}
	}
	return nil
}

func validateBuisnessLogic(product internal.Product) error {
	if product.Name == "" {
		return fmt.Errorf("%w: Name", internal.ErrFieldRequire)
	}
	if product.Quantity == 0 {
		return fmt.Errorf("%w: Quantity", internal.ErrFieldRequire)
	}
	if product.CodeValue == "" {
		return fmt.Errorf("%w: CodeValue", internal.ErrFieldRequire)
	}

	if product.Expiration == "" {
		return fmt.Errorf("%w: Expiration", internal.ErrFieldRequire)
	}

	regexpExpirationFormat := `^(0[1-9]|[1-2][0-9]|3[0-1])/(0[1-9]|1[0-2])/\d{4}$`
	expirationFormat, _ := regexp.MatchString(regexpExpirationFormat, product.Expiration)
	if !expirationFormat {
		return fmt.Errorf("%w: expiration", internal.ErrInvalidFormat)
	}

	if product.Price == 0.0 {
		return fmt.Errorf("%w: Price", internal.ErrFieldRequire)
	}
	return nil
}
