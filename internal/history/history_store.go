package history

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
)

// SaveDeploymentResult appends a deployment result to local JSONL history.
func SaveDeploymentResult(
	storageDirectory string,
	result *engine.DeploymentResult,
) error {
	if result == nil {
		return nil
	}

	if storageDirectory == "" {
		storageDirectory = ".lunar-deploy"
	}

	if err := os.MkdirAll(storageDirectory, 0755); err != nil {
		return err
	}

	historyPath := filepath.Join(
		storageDirectory,
		"history.jsonl",
	)

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		historyPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(string(resultJSON) + "\n"); err != nil {
		return err
	}

	return nil
}
