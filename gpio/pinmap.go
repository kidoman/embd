package gpio

type pinMap []*pinDesc

func (m pinMap) lookup(k interface{}) (*pinDesc, bool) {
	switch key := k.(type) {
	case int:
		for i := range m {
			if m[i].n == key {
				return m[i], true
			}
		}
	case string:
		for i := range m {
			for j := range m[i].ids {
				if m[i].ids[j] == key {
					return m[i], true
				}
			}
		}
	}

	return nil, false
}
