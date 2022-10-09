package heartbeat

func ChooseRandomDataServers(n int, excluded []string) []string {
	var candidates []string
	reverseExcluded := make(map[string]struct{})
	for i := range excluded {
		reverseExcluded[excluded[i]] = struct{}{}
	}
	alivable := getDataServers()
	for i := range alivable {
		if _, ok := reverseExcluded[alivable[i]]; !ok {
			candidates = append(candidates, alivable[i])
		}
	}
	return candidates
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