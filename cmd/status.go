/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/docker/go-units"
	"github.com/spf13/cobra"
	"github.com/wtran29/go-orchestrator/task"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status comman to list tasks",
	Long: `archon status command.

The status command allows a user to get the status of tasks from the Archon manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, _ := cmd.Flags().GetString("manager")
		url := fmt.Sprintf("http://%s/tasks", manager)
		resp, _ := http.Get(url)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var tasks []*task.Task
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tNAME\tCREATED\tSTATE\tCONTAINERNAME\tIMAGE\t")
		for _, task := range tasks {
			var start string
			if task.StartTime.IsZero() {
				start = fmt.Sprintf("%s ago", units.HumanDuration(time.Now().UTC().Sub(time.Now().UTC())))
			} else {
				start = fmt.Sprintf("%s ago", units.HumanDuration(time.Now().UTC().Sub(task.StartTime)))
			}
			state := task.State.String()[task.State]
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", task.ID, task.Name, start, state, task.Name, task.Image)

		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("manager", "m", "localhost:5555", "Manager to talk to")

}
