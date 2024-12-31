package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"github.com/invopop/yaml"
	"github.com/rs/zerolog/log"
	. "go-spring/domains/repository"
	"os"
	"os/exec"
)

type PlaybookExecutorStreamingOutput struct {
	ID      string
	Content string
}

type PlaybookExecutor struct {
	StdOutChannel chan PlaybookExecutorStreamingOutput
	StdErrChannel chan PlaybookExecutorStreamingOutput
}

func NewPlaybookExecutor() *PlaybookExecutor {
	return &PlaybookExecutor{
		StdOutChannel: make(chan PlaybookExecutorStreamingOutput),
		StdErrChannel: make(chan PlaybookExecutorStreamingOutput),
	}
}

func (executor *PlaybookExecutor) ExecutePlaybook(ctx context.Context, book *Playbook, inventory *Inventory) error {
	log.Info().Msgf("Writing and executing playbook: %v", book.ID)
	bookFile, err := executor.writePlaybook(book)
	if err != nil {
		log.Error().Msgf("Error writing playbook: %v", err)
		return err
	}

	inventoryFile, err := executor.writeInventory(inventory)
	if err != nil {
		log.Error().Msgf("Error writing inventory: %v", err)
		return err
	}

	var ansiblePlaybookRunTimeError error
	if err := executor.executePlaybook(ctx, book.ID, bookFile, inventoryFile); err != nil {
		log.Error().Msgf("Error executing playbook: %v", err)
		ansiblePlaybookRunTimeError = err
	}
	log.Info().Msgf("Deleting playbook: %v", book.ID)
	if err := executor.deletePlaybook(book); err != nil {
		log.Info().Msgf("Error deleting playbook: %v", err)
		return err
	}

	log.Info().Msgf("Deleting inventory: %v", inventory.ID)
	if err := executor.deleteInventory(inventory); err != nil {
		log.Info().Msgf("Error deleting inventory: %v", err)
		return err
	}

	return ansiblePlaybookRunTimeError
}

func (executor *PlaybookExecutor) executePlaybook(ctx context.Context, streamId string, bookFile string, inventoryFile string) error {
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		VerboseVV: true,
		Inventory: inventoryFile,
	}

	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
		playbook.WithPlaybooks(bookFile),
	)

	commands, _ := playbookCmd.Command()
	playbookEntryPoint, args := commands[0], commands[1:]

	log.Info().Msgf("Running command: %v", commands)

	return executor.runCommand(ctx, streamId, playbookEntryPoint, args...)
}

func (executor *PlaybookExecutor) runCommand(ctx context.Context, streamId string, command string, args ...string) error {
	_ = os.Setenv("ANSIBLE_FORCE_COLOR", "true")

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled")
	default:
		cmd := exec.Command(command, args...)

		// Get the command's stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return err
		}

		// Create readers to stream the output
		stdoutScanner := bufio.NewScanner(stdout)
		stderrScanner := bufio.NewScanner(stderr)

		// Stream stdout
		go func() {
			for stdoutScanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
					executor.StdOutChannel <- PlaybookExecutorStreamingOutput{
						ID:      streamId,
						Content: stdoutScanner.Text(),
					}
				}
			}
		}()

		// Stream stderr
		go func() {
			for stderrScanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
					executor.StdErrChannel <- PlaybookExecutorStreamingOutput{
						ID:      streamId,
						Content: stderrScanner.Text(),
					}
				}
			}
		}()

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			return err
		}

		log.Info().Msgf("Command %v executed successfully", command)
		return nil
	}

}

func (executor *PlaybookExecutor) getPlaybookLocation(book *Playbook) string {
	return fmt.Sprintf("/tmp/%s.yml", book.ID)
}

func (executor *PlaybookExecutor) getInventoryLocation(inventory *Inventory) string {
	return fmt.Sprintf("/tmp/%s.yml", inventory.ID)
}

func (executor *PlaybookExecutor) handleWrite(filePath string, content string) (string, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		return "", err
	}

	if _, err = file.WriteString(content); err != nil {
		return "", err
	}
	return filePath, err
}

func (executor *PlaybookExecutor) writeInventory(inventory *Inventory) (string, error) {
	inventoryFile := executor.getInventoryLocation(inventory)

	jsonMap := inventory.AsJson()
	_json, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}
	_yaml, err := yaml.JSONToYAML(_json)
	if err != nil {
		return "", err
	}

	return executor.handleWrite(inventoryFile, string(_yaml))
}

func (executor *PlaybookExecutor) writePlaybook(book *Playbook) (string, error) {
	bookFile := executor.getPlaybookLocation(book)
	jsonMap := book.AsJson()
	_json, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	_yaml, err := yaml.JSONToYAML(_json)
	if err != nil {
		return "", err
	}

	return executor.handleWrite(bookFile, string(_yaml))
}

func (executor *PlaybookExecutor) deletePlaybook(book *Playbook) error {
	bookFile := executor.getPlaybookLocation(book)
	if err := os.Remove(bookFile); err != nil {
		return err
	}
	return nil
}

func (executor *PlaybookExecutor) deleteInventory(inventory *Inventory) error {
	inventoryFile := executor.getInventoryLocation(inventory)
	if err := os.Remove(inventoryFile); err != nil {
		return err
	}
	return nil
}
