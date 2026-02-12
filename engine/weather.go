package engine

var conditions = []string{
	"Clear skies", "Partly cloudy", "Overcast", "Light rain",
	"Heavy rain", "Thunderstorm", "Fog / mist", "Snow (light)",
	"Snow (heavy)", "Hail", "Strong wind", "Extreme weather",
}

var temps = [7]string{
	"", "Bitterly cold", "Cold", "Cool", "Mild", "Warm", "Hot",
}

var winds = [7]string{
	"", "Dead calm", "Light breeze", "Breezy", "Windy", "Strong gusts", "Gale force",
}

func RandomWeather(rng *Randomizer) WeatherResult {
	return WeatherResult{
		Condition:   conditions[rng.Intn(len(conditions))],
		Temperature: temps[rng.RollD6()],
		Wind:        winds[rng.RollD6()],
	}
}
