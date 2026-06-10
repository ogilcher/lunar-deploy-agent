package history

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
)

// ReadDeploymentHistory loads deployment results from a JSONL history file.
func ReadDeploymentHistory(
	historyPath string,
) ([]engine.DeploymentResult, error) {
	file, err := os.Open(historyPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	results := []engine.DeploymentResult{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		var result engine.DeploymentResult

		if err := json.Unmarshal(line, &result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
