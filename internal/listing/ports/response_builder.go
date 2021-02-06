package ports

import "github.com/crounch-me/back/internal/listing/domain/lists"

type ResponseBuilder struct {
	builder  *ResponseBuilder
	response *ListResponse
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: &ListResponse{},
	}
}

func (b *ResponseBuilder) FromDomain(l *lists.List) *ResponseBuilder {
	products := make([]*ProductResponse, 0)
	for _, p := range l.Products() {
		product := &ProductResponse{
			UUID: p.UUID(),
		}
		products = append(products, product)
	}

	contributors := make([]*ContributorResponse, 0)
	for _, c := range l.Contributors() {
		contributor := &ContributorResponse{
			UUID: c.UUID(),
		}
		contributors = append(contributors, contributor)
	}

	b.response.UUID = l.UUID()
	b.response.Name = l.Name()
	b.response.CreationDate = l.CreationDate()
	b.response.Contributors = contributors
	b.response.Products = products

	return b
}

func (b *ResponseBuilder) Reset() {
	b.response = &ListResponse{}
}

func (b *ResponseBuilder) Build() *ListResponse {
	return b.response
}
