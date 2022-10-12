package heartbeat

import "math/rand"

func ChooseRandomDataServers(n int, excluded map[string]int) []string {
	var candidates []string
	alivable := getDataServers()
	for i := range alivable {
		if _, ok := excluded[alivable[i]]; !ok {
			candidates = append(candidates, alivable[i])
		}
	}
	length := len(candidates)
	if length < n {
		return candidates
	}
	p := rand.Perm(length)
	var ds []string
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return ds
}

func getDataServers() []string {
	var alivable []string
	rmu.RLock()
	for dataServer := range dataServers {
		alivable = append(alivable, dataServer)
	}
	rmu.RUnlock()
	return alivable
}