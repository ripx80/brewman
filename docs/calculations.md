# Calculation

```text
+ Wasser
+ Sudhausausbeute
+ Schüttung

AF = Ausbeutefaktor / get from plato table
AS = Sudausbeute / SudYield
AM = Ausschlagmenge/Würze / DecisiveSeasoning
S  = Schüttung
SG = Dichte / get from plato table
SW = %gew/ Grad Plato/Stammwürze in Prozent / OriginalWort

	Skalieren bei Änderung
	[x] Schüttung
	[x] Ausschlagwürze
	[x] Sudausbeute

	Schüttung: SumMalt

	Ausschlagwürze: DecisiveSeasoning AM
	Sudausbeute: 	SudYield

	add tablewriter for output plato table

	schüttung in prozent ausrechnen
	Gewitchsprozent ist die Stammwürze

	[x] Hopfen/Stopfhopfen
	[x] Malze
	[x] Wasser
	[x] Incredients

	1. Masse der Stammwürze berechnen
		- Stammwürze und Ausschlagsmenge
		- SG Specific Gravity (Extraktgehalt)
		- Lookup SW -> SG und VOL
		Masse der Würze = 22 Liter x 1057 : 1000 = 23,25 kg
	2. Extraktmenge berechnen
	3. Extraktanteil berechnen
	4. Gesamtmasse der Schüttung berechnen
	5. Malzmengen berechnen

https://zwieselbrau.wordpress.com/2017/04/26/berechnung-der-schuettung/
Formeln:
	Sudhausausbeute [%] = Stammwürze[°P] * m ( Würzemenge[kg] ) / m (Schüttung [kg] )
	Schüttung[kg] = Würzemenge[kg] * Stammwürze[°P] / Sudhausausbeute[%]
	SG(Massendichte SG [kg / L])  = ( 4,13 (kg/m3) * Stammwürze (°P)  ) / 1000 + d(Wasser bei T=20°C)

	m ( Extraktwürze [kg] ) = V ( Würzmenge [L]) * d ( Massendichte SG [kg/L] )
	m ( benötigtes Extrakt [kg] ) = m ( Extraktwürze [kg] ) * Stammwürze (°P) / 100
	m (Masse Schüttung [kg]) = m( benötigtes Extrakt [kg] )  / Sudhausausbeute[%]

	abweichungen: https://www.maischemalzundmehr.de/index.php?inhaltmitte=toolssudhausausbeute
	Sudhausausbeute [%] = (Volumen Ausschlagwürze [L] · Spezifische Dichte [kg/L] · Dezimalwert Stammwürze · Temperaturfaktor / Schüttung [kg] ) · 100
	Schüttung [kg] = (Volumen Ausschlagwürze [L] · Spezifische Dichte [kg/L] · Dezimalwert Stammwürze / Sudhausausbeute [%]) · 100
	Spezifische Dichte [kg/L] = (Extrakt [°P] / (258,6 - (Extrakt [°P] / 258,2) · 227,1))+1
```

