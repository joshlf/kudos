package config

func mapFromProbs(p []Problem) map[string]Problem {
	m := make(map[string]Problem)
	for _, pp := range p {
		m[pp.Code] = pp
	}
	return m
}
