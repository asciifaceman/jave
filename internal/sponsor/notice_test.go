package sponsor

import "testing"

func TestParseMode(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    Mode
		wantErr bool
	}{
		{name: "default empty", raw: "", want: ModeFull},
		{name: "full", raw: "full", want: ModeFull},
		{name: "full case-insensitive", raw: "FULL", want: ModeFull},
		{name: "redacted", raw: "redacted", want: ModeRedacted},
		{name: "off", raw: "off", want: ModeOff},
		{name: "invalid", raw: "wat", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMode(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseMode returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseMode(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestRenderLines(t *testing.T) {
	full := RenderLines(ModeFull, "test-seed")
	if len(full) != 2 {
		t.Fatalf("full lines = %d, want 2", len(full))
	}
	if full[0] == "" || full[1] == "" {
		t.Fatal("full lines should be non-empty")
	}

	redacted := RenderLines(ModeRedacted, "test-seed")
	if len(redacted) != 2 {
		t.Fatalf("redacted lines = %d, want 2", len(redacted))
	}
	if redacted[0] == full[0] && redacted[1] == full[1] {
		t.Fatal("redacted lines should differ from full lines")
	}

	off := RenderLines(ModeOff, "test-seed")
	if len(off) != 0 {
		t.Fatalf("off lines = %d, want 0", len(off))
	}
}
