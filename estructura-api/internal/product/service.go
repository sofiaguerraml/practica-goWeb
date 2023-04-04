package product

import (
	"errors"
	"practica-goWeb/estructura-api/internal/domain"
)

type Service interface {
	Save(product domain.Product) (domain.Product, error)
	GetMaxPrice(price float64) ([]domain.Product, error)
	GetId(id int) (domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
	Update(id int, p domain.Product) (domain.Product, error)
	Delete(id int) error
}

type service struct {
	r Repository
}

//var (
// 	ErrProductInvalid  = errors.New("invalid product")
// 	ErrProductInternal = errors.New("internal error")
// )

// controller
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Save(product domain.Product) (domain.Product, error) {
	product, err := s.r.Save(product)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) GetMaxPrice(price float64) ([]domain.Product, error) {
	list := s.r.GetMaxPrice(price)
	if len(list) == 0 {
		return []domain.Product{}, errors.New("no products found")
	}
	return list, nil
}

func (s *service) GetId(id int) (domain.Product, error) {
	p, err := s.r.GetId(id)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (s *service) GetAllProducts() ([]domain.Product, error) {
	list := s.r.GetAllProducts()
	return list, nil
}

func (s *service) Delete(id int) error {
	err := s.r.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un producto
func (s *service) Update(id int, u domain.Product) (domain.Product, error) {
	p, err := s.r.GetId(id)
	if err != nil {
		return domain.Product{}, err
	}
	if u.Name != "" {
		p.Name = u.Name
	}
	if u.Code_value != "" {
		p.Code_value = u.Code_value
	}
	if u.Expiration != "" {
		p.Expiration = u.Expiration
	}
	if u.Quantity > 0 {
		p.Quantity = u.Quantity
	}
	if u.Price > 0 {
		p.Price = u.Price
	}
	p, err = s.r.Update(id, p)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}
