package protobuf

func EntitiesToProtobuf[T any, P any](entities []T, toProtobuf func(T) P) []P {
	var pEntities []P
	for _, entity := range entities {
		pEntities = append(pEntities, toProtobuf(entity))
	}
	return pEntities
}
