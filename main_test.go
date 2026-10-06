package main

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		requests    int
		concurrency int
		wantError   bool
	}{
		{
			name:        "parâmetros válidos",
			url:         "http://example.com",
			requests:    100,
			concurrency: 10,
			wantError:   false,
		},
		{
			name:        "URL obrigatória",
			url:         "",
			requests:    100,
			concurrency: 10,
			wantError:   true,
		},
		{
			name:        "requests deve ser maior que zero",
			url:         "http://example.com",
			requests:    0,
			concurrency: 10,
			wantError:   true,
		},
		{
			name:        "concurrency deve ser maior que zero",
			url:         "http://example.com",
			requests:    100,
			concurrency: 0,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.url, tt.requests, tt.concurrency)

			if (err != nil) != tt.wantError {
				t.Fatalf("validate() erro = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}
