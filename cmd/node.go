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

	"github.com/spf13/cobra"
	"github.com/wtran29/go-orchestrator/node"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node command to list nodes.",
	Long: `cube node command.

The node command allows a user to get the information about the nodes in the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, _ := cmd.Flags().GetString("manager")
		url := fmt.Sprintf("http://%s/nodes", manager)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var nodes []*node.Node
		json.Unmarshal(body, &nodes)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "NAME\tMEMORY (MiB)\tDISK (GiB)\tROLE\tTASKS\t")
		for _, node := range nodes {
			fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%d\t\n", node.Name, node.Memory/1000, node.Disk/1000/1000/1000, node.Role, node.TaskCount)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringP("manager", "m", "localhost:5555", "Manager to talk to")

}
