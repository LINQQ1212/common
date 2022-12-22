package models

func (x *Product) Copy(p *Product) {
	x.ID = p.ID
	x.CateId = p.CateId
	x.DomainID = p.DomainID
	x.Name = p.Name
	x.Image = p.Image
	x.Link = p.Link
	x.Keywords = p.Keywords
	x.Brand = p.Brand
}
