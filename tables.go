package dbr

type tables []*table

func (t tables) Build() (string, error) {

	var query string

	lenT := len(t)

	for i, item := range t {
		value, err := item.Build()
		if err != nil {
			return "", err
		}

		query += value

		if i+1 < lenT {
			query += ", "
		}
	}

	return query, nil
}
