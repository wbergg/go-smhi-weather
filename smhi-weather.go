package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

// Parameter struct for SMHI API response
type Parameter struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

// TimeSeriesEntry represents a single forecast entry
type TimeSeriesEntry struct {
	ValidTime  string      `json:"validTime"`
	Parameters []Parameter `json:"parameters"`
}

// SMHIResponse represents the complete SMHI API response
type SMHIResponse struct {
	TimeSeries []TimeSeriesEntry `json:"timeSeries"`
}

func getSMHIForecast(lat, lon float64) (*SMHIResponse, error) {
	url := fmt.Sprintf("https://opendata-download-metfcst.smhi.se/api/category/pmp3g/version/2/geotype/point/lon/%.4f/lat/%.4f/data.json",
		lon, lat)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data SMHIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func getParameterValue(params []Parameter, name string) *float64 {
	for _, p := range params {
		if p.Name == name && len(p.Values) > 0 {
			return &p.Values[0]
		}
	}
	return nil
}

func getWeatherConditionFromSymbol(symbolCode int) string {
	switch symbolCode {
	case 1:
		return "clear_sky"
	case 2:
		return "nearly_clear_sky"
	case 3:
		return "variable_cloudiness"
	case 4:
		return "halfclear_sky"
	case 5:
		return "cloudy_sky"
	case 6:
		return "overcast"
	case 7:
		return "fog"
	case 8, 9, 10:
		return "rain_showers"
	case 11:
		return "thunderstorm"
	case 12, 13, 14:
		return "sleet_showers"
	case 15, 16, 17:
		return "snow_showers"
	case 18, 19, 20:
		return "rain"
	case 21:
		return "thunder"
	case 22, 23, 24:
		return "sleet"
	case 25, 26, 27:
		return "snowfall"
	default:
		return "clear_sky"
	}
}

func getWeatherIcon(condition string) []string {
	icons := map[string][]string{
		"clear_sky": {
			"    \\   /    ",
			"     .-.     ",
			"  ― (   ) ―  ",
			"     `-'     ",
			"    /   \\    ",
		},
		"nearly_clear_sky": {
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__).",
			"             ",
		},
		"variable_cloudiness": {
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__).",
			"             ",
		},
		"halfclear_sky": {
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__).",
			"             ",
		},
		"cloudy_sky": {
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
		},
		"overcast": {
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
		},
		"fog": {
			"             ",
			" _ - _ - _ - ",
			"  _ - _ - _  ",
			" _ - _ - _ - ",
			"             ",
		},
		"rain_showers": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ' ' ' '  ",
			"   ' ' ' '   ",
		},
		"thunderstorm": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ⚡ ⚡ ⚡   ",
			"  ‚'‚'‚'‚'   ",
		},
		"sleet_showers": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    * ' * '  ",
			"   * ' * '   ",
		},
		"snow_showers": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    *  *  *  ",
			"   *  *  *   ",
		},
		"rain": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"  ‚'‚'‚'‚'   ",
			"  ‚'‚'‚'‚'   ",
		},
		"thunder": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ⚡ ⚡ ⚡   ",
			"  ‚'‚'‚'‚'   ",
		},
		"sleet": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"  * ' * ' *  ",
			"  * ' * ' *  ",
		},
		"snowfall": {
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"   *  *  *   ",
			"  *  *  *  * ",
		},
	}
	if icon, ok := icons[condition]; ok {
		return icon
	}
	return icons["clear_sky"]
}

func getWindDirection(degrees float64) string {
	// Compass calculator
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((degrees + 11.25) / 22.5)
	return directions[index%16]
}

func calculateApparentTemp(temp, windSpeed, humidity float64) float64 {
	// For cold temperatures, calculate wind chill
	if temp <= 10.0 && windSpeed >= 1.3 {
		// Wind chill formula (metric)
		windKmh := windSpeed * 3.6 // Convert m/s to km/h
		return 13.12 + 0.6215*temp - 11.37*math.Pow(windKmh, 0.16) + 0.3965*temp*math.Pow(windKmh, 0.16)
	}

	// For hot temperatures, calculate heat index
	if temp >= 27.0 && humidity >= 40.0 {
		// Heat index formula (simplified)
		c1 := -8.78469475556
		c2 := 1.61139411
		c3 := 2.33854883889
		c4 := -0.14611605
		c5 := -0.012308094
		c6 := -0.0164248277778
		c7 := 0.002211732
		c8 := 0.00072546
		c9 := -0.000003582

		t2 := temp * temp
		h2 := humidity * humidity

		hi := c1 + c2*temp + c3*humidity + c4*temp*humidity +
			c5*t2 + c6*h2 + c7*t2*humidity + c8*temp*h2 + c9*t2*h2
		return hi
	}

	// For moderate temperatures, return actual temperature
	return temp
}

func main() {
	// Mälarhöjden, Stockholm coordinates
	const (
		stockholmLat = 59.3009642
		stockholmLon = 17.9557798
	)

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║            Weather for Mälarhöjden, Stockholm             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Fetch data
	data, err := getSMHIForecast(stockholmLat, stockholmLon)
	if err != nil {
		fmt.Printf("Error fetching weather: %v\n", err)
		os.Exit(1)
	}

	if len(data.TimeSeries) == 0 {
		fmt.Println("No forecast data available")
		os.Exit(1)
	}

	// Get current forecast (first entry from API)
	// API returns UTC time
	current := data.TimeSeries[0]
	validTime, _ := time.Parse(time.RFC3339, current.ValidTime)
	// Convert to local time
	validTime = validTime.Local()

	// Extract parameters
	temp := getParameterValue(current.Parameters, "t")
	windSpeed := getParameterValue(current.Parameters, "ws")
	windDir := getParameterValue(current.Parameters, "wd")
	humidity := getParameterValue(current.Parameters, "r")
	pressure := getParameterValue(current.Parameters, "msl")
	visibility := getParameterValue(current.Parameters, "vis")
	weatherSymbol := getParameterValue(current.Parameters, "Wsymb2")

	// Get  weather condition symbol
	condition := "clear_sky"
	if weatherSymbol != nil {
		condition = getWeatherConditionFromSymbol(int(*weatherSymbol))
	}

	// Get ASCII icon for print
	icon := getWeatherIcon(condition)

	// Print current weather with icon
	fmt.Printf("┌─────────────────────────────┬─────────────────────────────┐\n")
	fmt.Printf("│ %s │", validTime.Format("Mon 02 Jan 15:04"))
	fmt.Printf("%s│\n", strings.Repeat(" ", 28-len(validTime.Format("Mon 02 Jan 15:04"))))

	for i, line := range icon {
		if i == 0 {
			fmt.Printf("│%s", line)
			if temp != nil {
				tempStr := fmt.Sprintf("%.1f°C", *temp)
				fmt.Printf("│ Temperature:  %s", tempStr)
				fmt.Printf("%s│\n", strings.Repeat(" ", 14-len(tempStr)))
			} else {
				fmt.Printf("│ Temperature:  N/A           │\n")
			}
		} else if i == 1 {
			fmt.Printf("│%s", line)
			if windSpeed != nil && windDir != nil {
				windStr := fmt.Sprintf("%.1f m/s %s", *windSpeed, getWindDirection(*windDir))
				fmt.Printf("│ Wind:         %s", windStr)
				fmt.Printf("%s│\n", strings.Repeat(" ", 14-len(windStr)))
			} else {
				fmt.Printf("│ Wind:         N/A           │\n")
			}
		} else if i == 2 {
			fmt.Printf("│%s", line)
			if humidity != nil {
				humStr := fmt.Sprintf("%.0f%%", *humidity)
				fmt.Printf("│ Humidity:     %s", humStr)
				fmt.Printf("%s│\n", strings.Repeat(" ", 14-len(humStr)))
			} else {
				fmt.Printf("│ Humidity:     N/A           │\n")
			}
		} else if i == 3 {
			fmt.Printf("│%s", line)
			if pressure != nil {
				presStr := fmt.Sprintf("%.0f hPa", *pressure)
				fmt.Printf("│ Pressure:     %s", presStr)
				fmt.Printf("%s│\n", strings.Repeat(" ", 14-len(presStr)))
			} else {
				fmt.Printf("│ Pressure:     N/A           │\n")
			}
		} else if i == 4 {
			fmt.Printf("│%s", line)
			if visibility != nil {
				visStr := fmt.Sprintf("%.1f km", *visibility)
				fmt.Printf("│ Visibility:   %s", visStr)
				fmt.Printf("%s│\n", strings.Repeat(" ", 14-len(visStr)))
			} else {
				fmt.Printf("│ Visibility:   N/A           │\n")
			}
		}
	}
	fmt.Printf("└─────────────────────────────┴─────────────────────────────┘\n")

	// Print hourly forecast
	fmt.Println("\n┌────────────────────────────────────────────── Hourly Forecast ──────────────────────────────────────────┐")
	fmt.Println("│ Time  │  Temp │ Feels like │   Wind     │  Rain  │ Humidity │ Pressure │ Visibility │ Condition         │")
	fmt.Println("├───────┼───────┼────────────┼────────────┼────────┼──────────┼──────────┼────────────┼───────────────────┤")

	// Show next 12 hours
	for i := 0; i < 12 && i < len(data.TimeSeries); i++ {
		entry := data.TimeSeries[i]
		t, _ := time.Parse(time.RFC3339, entry.ValidTime)
		t = t.Local() // Convert to local time

		temp := getParameterValue(entry.Parameters, "t")
		wind := getParameterValue(entry.Parameters, "ws")
		windDir := getParameterValue(entry.Parameters, "wd")
		rain := getParameterValue(entry.Parameters, "pmean")
		humidity := getParameterValue(entry.Parameters, "r")
		pressure := getParameterValue(entry.Parameters, "msl")
		vis := getParameterValue(entry.Parameters, "vis")
		weatherSymbol := getParameterValue(entry.Parameters, "Wsymb2")

		// Get  weather condition symbol
		cond := "clear_sky"
		if weatherSymbol != nil {
			cond = getWeatherConditionFromSymbol(int(*weatherSymbol))
		}

		// Detailed condition text with descriptions
		// Note: Each emoji takes 2 visual columns, so padding accounts for visual width = 18 columns
		condText := map[string]string{
			"clear_sky":           "☀️  Clear sky     ",
			"nearly_clear_sky":    "🌤️  Nearly clear  ",
			"variable_cloudiness": "⛅ Variable clouds",
			"halfclear_sky":       "⛅ Half clear    ",
			"cloudy_sky":          "☁️  Cloudy        ",
			"overcast":            "☁️  Overcast      ",
			"fog":                 "🌫️  Fog           ",
			"rain_showers":        "🌦️  Rain showers  ",
			"thunderstorm":        "⛈️  Thunderstorm  ",
			"sleet_showers":       "🌨️  Sleet showers ",
			"snow_showers":        "🌨️  Snow showers  ",
			"rain":                "🌧️  Rain          ",
			"thunder":             "⛈️  Thunder       ",
			"sleet":               "🌨️  Sleet         ",
			"snowfall":            "❄️  Snowfall      ",
		}

		// Fixed column widths (characters between pipes, excluding the pipes themselves)
		const (
			timeWidth       = 5
			tempWidth       = 5
			feelsLikeWidth  = 10
			windWidth       = 10
			rainWidth       = 6
			humidityWidth   = 8
			pressureWidth   = 8
			visibilityWidth = 10
			conditionWidth  = 18
		)

		// Format each field to exact width
		tempStr := "N/A"
		if temp != nil {
			tempStr = fmt.Sprintf("%.1f°", *temp)
		}

		feelsLikeStr := "N/A"
		if temp != nil && wind != nil && humidity != nil {
			feelsLike := calculateApparentTemp(*temp, *wind, *humidity)
			feelsLikeStr = fmt.Sprintf("%.1f°", feelsLike)
		}

		windStr := "N/A"
		if wind != nil && windDir != nil {
			windStr = fmt.Sprintf("%.1fm/s %s", *wind, getWindDirection(*windDir))
		} else if wind != nil {
			windStr = fmt.Sprintf("%.1fm/s", *wind)
		}

		rainStr := "N/A"
		if rain != nil {
			rainStr = fmt.Sprintf("%.1fmm", *rain)
		}

		humidityStr := "N/A"
		if humidity != nil {
			humidityStr = fmt.Sprintf("%.0f%%", *humidity)
		}

		pressureStr := "N/A"
		if pressure != nil {
			pressureStr = fmt.Sprintf("%.0fhPa", *pressure)
		}

		visStr := "N/A"
		if vis != nil {
			visStr = fmt.Sprintf("%.1fkm", *vis)
		}

		fmt.Printf("│ %*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %s │\n",
			timeWidth, t.Format("15:04"),
			tempWidth, tempStr,
			feelsLikeWidth, feelsLikeStr,
			windWidth, windStr,
			rainWidth, rainStr,
			humidityWidth, humidityStr,
			pressureWidth, pressureStr,
			visibilityWidth, visStr,
			condText[cond])
	}
	fmt.Println("└───────┴───────┴────────────┴────────────┴────────┴──────────┴──────────┴────────────┴───────────────────┘")
}
