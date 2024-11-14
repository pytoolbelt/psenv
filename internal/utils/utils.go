package utils

type Parameters struct {
	ToAdd    map[string]string
	ToUpdate map[string]string
	ToDelete []string
}

func MergeLocalAndRemoteParams(localParams, remoteParams map[string]string) *Parameters {
	var toMerge = &Parameters{
		ToAdd:    make(map[string]string),
		ToUpdate: make(map[string]string),
		ToDelete: make([]string, 0),
	}

	for localKey, localValue := range localParams {
		remoteValue, remoteExists := remoteParams[localKey]

		// if the local param does not exist in the remote params, then it must be added
		if !remoteExists {
			toMerge.ToAdd[localKey] = localValue
			continue
		}

		// if the local param exists in the remote params, but the value is different, then it must be updated
		if localValue != remoteValue {
			toMerge.ToUpdate[localKey] = localValue
			continue
		}
	}
	// if the remote param does not exist in the local params, then it must be deleted
	for k := range remoteParams {
		_, exists := localParams[k]
		if !exists {
			toMerge.ToDelete = append(toMerge.ToDelete, k)
		}
	}

	return toMerge
}
