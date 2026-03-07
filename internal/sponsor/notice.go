package sponsor

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// Mode controls sponsor notice rendering behavior.
type Mode string

const (
	ModeFull     Mode = "full"
	ModeRedacted Mode = "redacted"
	ModeOff      Mode = "off"
)

// ParseMode parses notice mode text into a stable enum.
func ParseMode(raw string) (Mode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", string(ModeFull):
		return ModeFull, nil
	case string(ModeRedacted):
		return ModeRedacted, nil
	case string(ModeOff):
		return ModeOff, nil
	default:
		return "", fmt.Errorf("invalid sponsor notice mode %q (expected full|redacted|off)", raw)
	}
}

// Corporate sponsors for Orcal ecosystem sustainability.
var corporateSponsors = []string{
	"Orcal Cloud",
	"MegaDyne Platform Group",
	"Venture-Adjacent Learning Collective",
	"Synergy Dynamics International",
	"Enterprise Momentum Partners",
	"Global Alignment Solutions",
	"Strategic Initiative Consortium",
	"Paradigm Shift Technologies",
	"Holistic Value Architects",
	"Agile Transformation Guild",
	"Cloud Native Excellence Foundation",
	"Digital Acceleration Partners",
	"Innovation Enablement Systems",
	"Cross-Functional Optimization Lab",
	"Sustainable Growth Ventures",
}

// RenderLines returns deterministic sponsor notice lines for the selected mode.
// The seed parameter ensures deterministic rotation while providing variety.
func RenderLines(mode Mode, seed string) []string {
	switch mode {
	case ModeOff:
		return nil
	case ModeRedacted:
		return []string{
			"javec: ecosystem notice: This build is supported by approved sustainability partners.",
			"javec: ecosystem notice: sponsor roster redacted under Policy-7.4 transparency minimization.",
		}
	default:
		// Select 3 sponsors deterministically based on seed
		sponsors := selectSponsors(seed, 3)
		roster := strings.Join(sponsors, ", ")
		return []string{
			"javec: ecosystem notice: This build is supported by Orcal-approved sustainability partners.",
			"javec: ecosystem notice: strategic partner roster: " + roster + ".",
		}
	}
}

// selectSponsors returns n sponsors deterministically selected based on seed.
func selectSponsors(seed string, n int) []string {
	if n <= 0 || n > len(corporateSponsors) {
		n = 3
	}

	// Hash the seed to get a deterministic starting index
	h := fnv.New32a()
	h.Write([]byte(seed))
	offset := int(h.Sum32()) % len(corporateSponsors)

	// Select n distinct sponsors starting from offset
	result := make([]string, 0, n)
	seen := make(map[int]bool)

	for i := 0; len(result) < n && i < len(corporateSponsors); i++ {
		idx := (offset + i) % len(corporateSponsors)
		if !seen[idx] {
			result = append(result, corporateSponsors[idx])
			seen[idx] = true
		}
	}

	return result
}
