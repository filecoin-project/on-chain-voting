package proposal

import (
	"bufio"
	"fil-vote/model"
	"fil-vote/service"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

// ListProposalsCmd returns a command that lists proposals with pagination.
func ListProposalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List proposals with pagination",
		Run: func(cmd *cobra.Command, args []string) {
			page := model.Page
			pageSize := model.PageSize
			showTable := true // Flag to control whether to display the table

			// Create a scanner for user input
			scanner := bufio.NewScanner(os.Stdin)

			// Continuous loop for displaying proposals and handling user commands
			for {
				proposals, err := service.GetProposalList(page, pageSize)
				if err != nil {
					fmt.Println("Error fetching proposals:", err)
					return
				}

				// If we're on the last page, don't show the table
				if showTable {
					displayProposals(proposals.Data.List)
				}

				// Display user interaction options
				displayMenu(page)

				// Read user input
				if scanner.Scan() {
					input := strings.TrimSpace(scanner.Text())

					// Process user input for pagination or selecting a proposal
					if exit := handleUserInput(input, &page, proposals.Data.Total, scanner, &showTable); exit {
						break // Exit the loop if 'q' is pressed
					}
				}
			}
		},
	}
	return cmd
}

// handleUserInput processes user commands for pagination or displaying proposal details.
func handleUserInput(input string, page *int, totalProposals int, scanner *bufio.Scanner, showTable *bool) bool {
	switch input {
	case "n":
		// Calculate if there is another page
		if *page*model.PageSize < totalProposals {
			*page++ // Next page
			*showTable = true
		} else {
			fmt.Println("End of proposals.")
			*showTable = false // Don't show the table if we're at the last page
		}
	case "p":
		if *page > 1 {
			*page-- // Previous page
			*showTable = true
		} else {
			fmt.Println("You are already on the first page.")
			*showTable = false // Don't show the table if we're at the first page
		}
	case "q":
		fmt.Println("Exiting to list...")
		return true // Indicate to break out of the loop
	default:
		// Handle if a Proposal ID is entered
		processProposalID(input, scanner)
	}
	return false // Keep the loop running
}

// processProposalID handles the user's input when a proposal ID is provided.
func processProposalID(input string, scanner *bufio.Scanner) {
	proposalID, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid Proposal ID. Please try again.")
		return
	}

	// Fetch detailed proposal content by ID
	proposal, err := service.GetProposalByID(int64(proposalID))
	if err != nil {
		fmt.Println("Error fetching proposal details:", err)
		return
	}
	if proposal.ProposalId == 0 {
		fmt.Println("Invalid Proposal ID. Please try again.")
		return
	}

	// Print detailed proposal content
	displayProposalDetails(proposal)

	// Wait for user input to return or continue
	for {
		fmt.Println("\nEnter 'q' to return to the proposal list or any other key to continue.")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "q" {
			return // If 'q' is pressed, return to the list
		} else {
			fmt.Println("\nEnter 'q' to return to the proposal list or any other key to continue.")
		}
	}
}

// displayProposals prints the list of proposals in a formatted table.
func displayProposals(proposals []model.Proposal) {
	table := createTable([]string{"Proposal ID", "Creator", "Title", "Status"})

	for _, proposal := range proposals {
		status := proposal.Status.String()
		table.Append([]string{
			fmt.Sprintf("%d", proposal.ProposalId),
			proposal.Address,
			proposal.Title,
			status,
		})
	}

	// Render the table
	table.Render()
}

// displayProposalDetails prints the detailed content of a proposal.
func displayProposalDetails(proposal model.Proposal) {
	proposalTable := createTable([]string{"Field", "Details"})

	proposalTable.Append([]string{"Proposal ID", fmt.Sprintf("%d", proposal.ProposalId)})
	proposalTable.Append([]string{"Creator", proposal.Address})
	proposalTable.Append([]string{"Title", proposal.Title})
	proposalTable.Append([]string{"Content", fmt.Sprintf("%-80s", proposal.Content)})
	proposalTable.Append([]string{"Start Time", time.Unix(proposal.StartTime, 0).Format("2006-01-02 15:04")})
	proposalTable.Append([]string{"End Time", time.Unix(proposal.EndTime, 0).Format("2006-01-02 15:04")})
	proposalTable.Append([]string{"Snapshot Block Height", fmt.Sprintf("%d", proposal.SnapshotInfo.SnapshotHeight)})
	proposalTable.Append([]string{"Status", proposal.Status.String()})
	if proposal.Status == model.ProposalStatusCompleted {
		proposalTable.Append([]string{"Vote Approve", fmt.Sprintf("%.2f%%", proposal.VotePercentage.Approve)})
		proposalTable.Append([]string{"Vote Reject", fmt.Sprintf("%.2f%%", proposal.VotePercentage.Reject)})
	}

	// Render the proposal content table
	proposalTable.Render()
}

// displayMenu shows the user interaction options.
func displayMenu(page int) {
	fmt.Println("\nEnter a command:")
	fmt.Println("n - Next page")
	fmt.Println("p - Previous page")
	fmt.Println("q - Quit to list")
	fmt.Println("Enter Proposal ID to view details:")
	fmt.Printf("You are currently on page %d.\n", page)
}

// createTable creates a new table with specified headers.
func createTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(true)
	table.SetColumnSeparator("|")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	return table
}
