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

func (res *Vec) Clamp(s *Vec) {
	for i := range *res {
		if (*res)[i] > (*s)[i]/2 {
			(*res)[i] = (*s)[i] / 2
		}
		if (*res)[i] < -(*s)[i]/2 {
			(*res)[i] = -(*s)[i] / 2
		}
	}
}

func Dot(a, b *Vec) float64 {
	return (*a)[0]*(*b)[0] + (*a)[1]*(*b)[1] + (*a)[2]*(*b)[2]
}

func (v *Vec) Nrm2Sq() float64 {
	return Dot(v, v)
}

func (res *Vec) Scale(alpha float64, v *Vec) *Vec {
	(*res)[0] = alpha * (*v)[0]
	(*res)[1] = alpha * (*v)[1]
	(*res)[2] = alpha * (*v)[2]
	return res
}
