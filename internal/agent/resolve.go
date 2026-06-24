package agent

// Resolve returns the agent adapters for the given names. If names is empty,
// it falls back to auto-detected agents, then to the "claude" built-in.
func Resolve(names []string) []Adapter {
	if len(names) > 0 {
		var result []Adapter
		for _, n := range names {
			if a, ok := Find(n); ok {
				result = append(result, a)
			}
		}
		return result
	}

	detected := Detect()
	if len(detected) > 0 {
		return detected
	}

	if a, ok := Find("claude"); ok {
		return []Adapter{a}
	}
	return nil
}
