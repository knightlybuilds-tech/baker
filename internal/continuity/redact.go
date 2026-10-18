package continuity

import "regexp"

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`),
	regexp.MustCompile(`\b(sk|pk|rk)[-_](live|test|proj|ant)?[-_]?[A-Za-z0-9_-]{16,}`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
	regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}`),
	regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
	regexp.MustCompile(`\bxox[abprs]-[A-Za-z0-9-]{10,}`),
	regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{30,}`),
	regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`), // JWT
	regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/-]{16,}=*`),
}

// key=value / key: value where the key looks secret.
var secretAssign = regexp.MustCompile(`(?i)\b([A-Z0-9_]*(api[_-]?key|secret|token|password|passwd|pwd)[A-Z0-9_]*)(\s*[:=]\s*)("[^"]*"|'[^']*'|\S+)`)

// Redact replaces anything that looks like a credential with [REDACTED].
func Redact(s string) string {
	for _, re := range secretPatterns {
		s = re.ReplaceAllString(s, "[REDACTED]")
	}
	return secretAssign.ReplaceAllString(s, "$1$3[REDACTED]")
}
