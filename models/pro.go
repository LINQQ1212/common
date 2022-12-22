package models

func (x *Product) Copy(p *Product) {
	x.ID = p.ID
	x.CateId = p.CateId
	x.DomainID = p.DomainID
	x.Pid = p.Pid
	x.Model = p.Model
	x.Image = p.Image
	x.Images = p.Images
	x.Price = p.Price
	x.Specials = p.Specials
	x.Name = p.Name
	x.Description = p.Description
	x.MTitle = p.MTitle
	x.MKeywords = p.MKeywords
	x.MDescription = p.MDescription
	x.Link = p.Link
	x.Keywords = p.Keywords
	x.CPath = p.CPath
	x.GoogleImgs = p.GoogleImgs
	x.YahooDesc = p.YahooDesc
	x.Categories = p.Categories
	x.Brand = p.Brand
	x.BingDesc = p.BingDesc
	x.Youtube = p.Youtube
}
