package envgen

func contains[e comparable](s []e, v e) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}

	return false
}
