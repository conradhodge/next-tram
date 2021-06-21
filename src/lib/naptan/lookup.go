package naptan

// LookupStop returns the National Public Transport Access Node (NaPTAN) code for the stop name given
func LookupStop(stopName string) string {
	stops := map[string]string{
		"Malin Bridge": "37090179",
		"Hillsborough": "37090177",
	}
	return stops[stopName]
}
