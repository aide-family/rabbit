package bo

func NewPageRequestBo[T any](page int32, pageSize int32) *PageRequestBo[T] {
	return &PageRequestBo[T]{
		Page:     page,
		PageSize: pageSize,
	}
}

type PageRequestBo[T any] struct {
	Page     int32
	PageSize int32
}

func (p *PageRequestBo[T]) NewPageResponseBo(data []T, total int64) *PageResponseBo[T] {
	return &PageResponseBo[T]{
		items:         data,
		total:         total,
		PageRequestBo: p,
	}
}

type PageResponseBo[T any] struct {
	items []T
	total int64
	*PageRequestBo[T]
}

func (p *PageResponseBo[T]) GetItems() []T {
	return p.items
}

func (p *PageResponseBo[T]) GetTotal() int64 {
	return p.total
}

func (p *PageResponseBo[T]) GetPage() int32 {
	return p.Page
}

func (p *PageResponseBo[T]) GetPageSize() int32 {
	return p.PageSize
}
