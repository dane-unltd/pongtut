package main

type Vec [3]float64

func (res *Vec) Add(a, b *Vec) *Vec {
	(*res)[0] = (*a)[0] + (*b)[0]
	(*res)[1] = (*a)[1] + (*b)[1]
	(*res)[2] = (*a)[2] + (*b)[2]
	return res
}

func (res *Vec) Sub(a, b *Vec) *Vec {
	(*res)[0] = (*a)[0] - (*b)[0]
	(*res)[1] = (*a)[1] - (*b)[1]
	(*res)[2] = (*a)[2] - (*b)[2]
	return res
}

func (a *Vec) Equals(b *Vec) bool {
	for i := range *a {
		if (*a)[i] != (*b)[i] {
			return false
		}
	}
	return true
}
