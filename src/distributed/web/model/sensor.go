package model

// Sensor describes a sensor
type Sensor struct {
	Name         string  `json:"name"`
	SerialNo     string  `json:"serial_no"`
	UnitType     string  `json:"unit_type"`
	MinSafeValue float64 `json:"min_safe_value"`
	MaxSafeValue float64 `json:"max_safe_value"`
}

// GetSensorByName queries for the sensor with the given name
func GetSensorByName(name string) (Sensor, error) {
	q := `SELECT name, serial_no, unit_type, min_saf_value, max_safe_value
					FROM sensor
					WHERE name = $1`
	result := Sensor{}
	row := db.QueryRow(q, name)
	err := row.Scan(&result.Name, &result.SerialNo, &result.UnitType, &result.MinSafeValue, &result.MaxSafeValue)
	return result, err

}
