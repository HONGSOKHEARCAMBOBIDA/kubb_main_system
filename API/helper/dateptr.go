package helper

func FormatDatePtr(date *string) *string {
	if date == nil {
		return nil
	}

	formatted := FormatDate(*date)
	return &formatted
}
