package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	if got != "hello" {
		t.Errorf("Hello() = %q, want \"hello\"", got)
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		// Gültige Eingaben
		{name: "1h30m", input: "1h30m", want: 5400},
		{name: "45s", input: "45s", want: 45},
		{name: "2d", input: "2d", want: 172800},
		{name: "1d2h3m4s", input: "1d2h3m4s", want: 93784},
		{name: "0s", input: "0s", want: 0},

		// Whitespace-Regeln — gültig
		{name: "Whitespace Leading/Trailing", input: " 1h30m ", want: 5400},
		{name: "Whitespace zwischen Einheit und Zahl", input: "1h 30m", want: 5400},
		{name: "Whitespace um alles", input: " 2d ", want: 172800},
		{name: "Whitespace überall außer inner-Token", input: " 1d 2h 3m 4s ", want: 93784},

		// Ungültige Eingaben
		{name: "leer", input: "", wantErr: true},
		{name: "unbekannte Einheit", input: "1x", wantErr: true},
		{name: "negative Zahl", input: "-1h", wantErr: true},
		{name: "doppelte Einheit h", input: "1h2h", wantErr: true},
		{name: "doppelte Einheit m", input: "5m30m", wantErr: true},

		// Whitespace-Regeln — ungültig
		{name: "Lücke Zahl zu Einheit", input: "1 h", wantErr: true},
		{name: "Lücke innerhalb Zahl", input: "1 0m", wantErr: true},
		{name: "doppelte Einheit mit Lücke", input: "1h 2h", wantErr: true},
		{name: "nur Leerzeichen", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDuration(%q) expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDuration(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
