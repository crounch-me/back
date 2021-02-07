package ports

import "github.com/crounch-me/back/internal/products/domain/products"

type ResponseBuilder struct {
	builder  *ResponseBuilder
	response *ProductResponse
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: &ProductResponse{},
	}
}

func (b *ResponseBuilder) FromDomain(p *products.Product) *ResponseBuilder {
	b.response.UUID = p.UUID()
	b.response.Name = p.Name()

	if p.HasCategory() {
		b.response.Category = &CategoryResponse{
			UUID: p.CategoryUUID(),
			Name: p.CategoryName(),
		}
	}

	return b
}

func (b *ResponseBuilder) Reset() {
	b.response = &ProductResponse{}
}

func (b *ResponseBuilder) Build() *ProductResponse {
	return b.response
}
