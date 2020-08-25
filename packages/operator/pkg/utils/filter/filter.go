package filter

type ListOptions struct {
	Filter interface{}
	Page   *int
	Size   *int
}

type ListOption func(*ListOptions)

func ListFilter(filter interface{}) ListOption {
	return func(args *ListOptions) {
		args.Filter = filter
	}
}

func Page(page int) ListOption {
	return func(args *ListOptions) {
		args.Page = &page
	}
}

func Size(size int) ListOption {
	return func(args *ListOptions) {
		args.Size = &size
	}
}
