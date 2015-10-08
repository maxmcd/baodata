package baodata

import "strconv"

func mockDelete(pd PathData) {
	resource := mockData[pd.ResourceName]
	idString := strconv.Itoa(pd.Id)
	var deleteIndex int
	for i, item := range resource {
		if item["id"] == idString {
			deleteIndex = i
		}
	}

	resource = append(
		resource[:deleteIndex],
		resource[deleteIndex+1:]...,
	)

	// completely ignore possible race conditions
	// with mock data
	mockData[pd.ResourceName] = resource
}

func mockGet(pd PathData) []Data {
	resource := mockData[pd.ResourceName]
	idString := strconv.Itoa(pd.Id)
	if pd.Id == 0 {
		return resource
	}
	for _, item := range resource {
		if item["id"] == idString {
			return []Data{item}
		}
	}
	return []Data{}
}

func mockPut(pd PathData, data Data) []Data {
	resource := mockData[pd.ResourceName]

	if pd.Id == 0 {
		newId := len(resource) + 1
		idString := strconv.Itoa(newId)
		data["id"] = idString
		resource = append(resource, data)
	} else {
		idString := strconv.Itoa(pd.Id)
		for i, item := range resource {
			if item["id"] == idString {
				for key, value := range data {
					resource[i][key] = value
				}
				data = resource[i]
				break
			}
		}
	}

	// completely ignore possible race conditions
	// with mock data
	mockData[pd.ResourceName] = resource

	return []Data{data}
}
