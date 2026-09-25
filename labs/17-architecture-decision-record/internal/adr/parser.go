package adr

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	titleRegex      = regexp.MustCompile(`^#\s+(\d+)\.\s+(.+)$`)
	statusRegex     = regexp.MustCompile(`(?i)^Status:\s*([A-Za-z]+)(?:\s+by\s+(\d+))?`)
	supersedesRegex = regexp.MustCompile(`(?i)^Supersedes:\s*(\d+)`)
)

func Parse(content string) (*Record, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	record := &Record{
		Content: content,
	}

	foundTitle := false
	foundStatus := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !foundTitle {
			matches := titleRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				id, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, fmt.Errorf("invalid record id: %w", err)
				}
				record.ID = id
				record.Title = matches[2]
				foundTitle = true
				continue
			}
		}

		if strings.HasPrefix(strings.ToLower(line), "status:") {
			matches := statusRegex.FindStringSubmatch(line)
			if len(matches) >= 2 {
				statusStr := strings.Title(strings.ToLower(matches[1]))
				record.Status = Status(statusStr)
				foundStatus = true

				if len(matches) == 3 && matches[2] != "" {
					supersededBy, err := strconv.Atoi(matches[2])
					if err == nil {
						record.SupersededBy = supersededBy
					}
				}
			}
		}

		if strings.HasPrefix(strings.ToLower(line), "supersedes:") {
			matches := supersedesRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				supersedes, err := strconv.Atoi(matches[1])
				if err == nil {
					record.Supersedes = supersedes
				}
			}
		}
	}

	if !foundTitle {
		return nil, fmt.Errorf("title not found or malformed")
	}

	if !foundStatus {
		return nil, fmt.Errorf("status not found")
	}

	if !record.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", record.Status)
	}

	return record, nil
}
