## SMHI API weather collector

Tool I used to collect and present weather from SMHI's API with wttr.in as inspiration.

No idea if this is of any use to anyone, SMHI is the Swedish Metrological Weather Institute. I guess it's just for me trying to learn GO and how to fetch and presenting data from an API.

## Usage

Change the location in the code to your location:

```
	// Mälarhöjden, Stockholm coordinates
	const (
		stockholmLat = 59.3009642
		stockholmLon = 17.9557798
	)
```

## Running

Run the program and you will get something like this back:

```

╔═══════════════════════════════════════════════════════════╗
║            Weather for Mälarhöjden, Stockholm             ║
╚═══════════════════════════════════════════════════════════╝

┌───────────────────────────────────────────────────────────┐
│ Thu 04 Dec 10:00                                          │
│             │ Temperature:  4.3°C                         │
│     .--.    │ Wind:         2.2 m/s S                     │
│  .-(    ).  │ Humidity:     91%                           │
│ (___.__)__) │ Pressure:     1015 hPa                      │
│             │ Visibility:   10.3 km                       │
└───────────────────────────────────────────────────────────┘

┌───────────────────────────────────────── Hourly Forecast ──────────────────────────────────────────────┐
│ Time  │ Temp  │ Feels like │ Wind        │ Rain   │ Humidity │ Pressure  │ Visibility  │ Condition     │
├───────┼───────┼────────────┼─────────────┼────────┼──────────┼───────────┼─────────────┼───────────────┤
│ 10:00 │  4.3° │     2.3°   │  2.2m/s S   │  0.0mm │     91%  │   1015hPa │    10.3km   │ ☁️  Overcast   │
│ 11:00 │  4.3° │     2.4°   │  2.1m/s S   │  0.0mm │     90%  │   1014hPa │    11.0km   │ ☁️  Overcast   │
│ 12:00 │  4.3° │     2.5°   │  2.0m/s SSE │  0.0mm │     92%  │   1014hPa │     9.7km   │ ☁️  Overcast   │
│ 13:00 │  4.4° │     2.8°   │  1.9m/s SSE │  0.0mm │     93%  │   1014hPa │     8.6km   │ ☁️  Overcast   │
│ 14:00 │  4.4° │     2.9°   │  1.8m/s SSE │  0.0mm │     94%  │   1014hPa │     7.9km   │ ☁️  Overcast   │
│ 15:00 │  4.5° │     2.8°   │  2.0m/s SSE │  0.0mm │     94%  │   1014hPa │    10.0km   │ ☁️  Overcast   │
│ 16:00 │  4.7° │     2.7°   │  2.3m/s S   │  0.0mm │     93%  │   1014hPa │     8.2km   │ ☁️  Overcast   │
│ 17:00 │  4.7° │     2.7°   │  2.3m/s SSE │  0.0mm │     93%  │   1015hPa │     8.7km   │ ☁️  Overcast   │
│ 18:00 │  4.6° │     2.4°   │  2.5m/s S   │  0.0mm │     93%  │   1015hPa │     8.9km   │ ☁️  Overcast   │
│ 19:00 │  4.3° │     2.6°   │  1.9m/s SSE │  0.0mm │     93%  │   1015hPa │     8.9km   │ ☁️  Overcast   │
│ 20:00 │  4.2° │     2.8°   │  1.7m/s SSE │  0.0mm │     95%  │   1015hPa │    10.0km   │ ☁️  Overcast   │
│ 21:00 │  4.1° │     2.8°   │  1.6m/s SSE │  0.0mm │     95%  │   1016hPa │    10.0km   │ ☁️  Overcast   │
└───────┴───────┴────────────┴─────────────┴────────┴──────────┴───────────┴─────────────┴───────────────┘

```