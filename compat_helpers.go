package forge

func compatListOptions(page, pageSize, limit int) (ListOptions, bool) {
	if page == 0 && pageSize == 0 && limit == 0 {
		return ListOptions{}, false
	}

	opts := ListOptions{
		Page:     page,
		PageSize: pageSize,
		Limit:    limit,
	}
	if opts.Page < 1 {
		opts.Page = 1
	}
	return opts, true
}
