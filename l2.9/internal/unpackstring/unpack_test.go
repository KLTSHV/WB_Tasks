package unpackstring

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple",
			input: "a4bc2d5e",
			want:  "aaaabccddddde",
		},
		{
			name:  "no digits",
			input: "abcd",
			want:  "abcd",
		},
		{
			name:    "only digits",
			input:   "45",
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "escaped digits separately",
			input: `qwe\4\5`,
			want:  "qwe45",
		},
		{
			name:  "escaped digit then multiplier",
			input: `qwe\45`,
			want:  "qwe44444",
		},
		{
			name:    "starts with digit",
			input:   "3abc",
			wantErr: true,
		},
		{
			name:    "dangling escape",
			input:   `qwe\`,
			wantErr: true,
		},
		{
			name:  "zero multiplier removes previous rune",
			input: "a0b",
			want:  "b",
		},
		{
			name:  "unicode rune",
			input: "Я4",
			want:  "ЯЯЯЯ",
		},
		{
			name:  "escape backslash itself then multiply",
			input: `qwe\\5`,
			want:  `qwe\\\\\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("Unpack() = %q, want %q", got, tt.want)
			}
		})
	}
}
