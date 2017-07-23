package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Set with go build -X
var version = "undefined"

// Struct to hold the contents of the command as well as the most recent
// timestamp and the index it was found in the original bash_history file
type command struct {
	cmd           string
	ts, origIndex int
}

// Sorts the array of commands in place.
// Sorts first by most recent usage (.ts), next by index it was found in the file.
func sortCommands(commands *[]command) {
	sort.SliceStable(*commands, func(i, j int) bool {
		if (*commands)[i].ts < (*commands)[j].ts {
			return true
		} else if (*commands)[i].ts > (*commands)[j].ts {
			return false
		} else {
			return (*commands)[i].origIndex < (*commands)[j].origIndex
		}
	})
}

// Deduplicates the array of commands.
// Goes through the (sorted) array of commands in reverse order. When duplicate
// commands are encountered, it keeps the usage with the most recent timestamp.
// https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/db6qa2e/
func dedupCommands(commands []command) []command {

	seen := make(map[string]struct{}, len(commands))
	j := len(commands) - 1

	for i := j; i >= 0; i-- {
		command := commands[i]
		if _, ok := seen[command.cmd]; ok {
			continue
		}
		seen[command.cmd] = struct{}{}
		commands[j] = command
		j--
	}
	return commands[j+1:]
}

// Split all the string of commmands into []command.
func makeCommandArr(lines []string) []command {

	var cmdFragment string
	var commands []command
	var startidx, ts int
	var err error

	// Regex for a line that is a timestamp
	re := regexp.MustCompile(`^#\d+$`)

	for idx, line := range lines {

		if re.MatchString(line) {
			// Line is a timestamp

			if len(cmdFragment) != 0 {
				commands = append(commands, command{cmd: cmdFragment, ts: ts, origIndex: startidx})
				cmdFragment = ""
			}

			ts, err = strconv.Atoi(line[1:])
			if err != nil {
				log.Fatal(err)
			}
			startidx = idx

		} else {
			// Line is not a timestamp, is either a command or a continuation
			// of a command

			// Will not yet containing anything on the first run through
			if len(cmdFragment) != 0 {
				cmdFragment = cmdFragment + "\n" + line
			} else {
				cmdFragment = line
			}
		}
	}

	// Append the last frament
	commands = append(commands, command{cmd: cmdFragment, ts: ts, origIndex: startidx})
	return commands
}

// Turn a []command back to a string for printing / writing
func commandsToString(commands []command) string {
	s := make([]string, len(commands)*2)
	for idx, command := range commands {
		s[idx*2] = "#" + strconv.Itoa(command.ts)
		s[idx*2+1] = command.cmd
	}
	return strings.Join(s, "\n")
}

func printVersion() {
	fmt.Println("fixbashhistory version:", version)
	os.Exit(0)
}

func main() {

	infilePath := flag.String("history-file", "", "Path to your .bash_history")
	outfilePath := flag.String("outfile", "", "Output file (defaults to stdout)")
	showVersion := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *showVersion {
		printVersion()
	}

	if *infilePath == "" {
		log.Fatal("You have to specify an input file.")
	}

	infile, err := os.Open(*infilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	var lines []string
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Make the string into an array of individual commands.
	commands := makeCommandArr(lines)

	// Sort so that most recent usage is last.
	sortCommands(&commands)

	// Uses a map to filter so the most recent (last) command is preserved.
	// Builds an array from the map which will unfortunately lose its sort order due to how maps work.
	commandsNoDupes := dedupCommands(commands)

	// Fix the sorting lost in deduplication
	sortCommands(&commandsNoDupes)

	var outfile *os.File

	// Default to printing to stdout if an outfile wasn't specified.
	if *outfilePath == "" {
		outfile = os.Stdout
	} else {
		outfile, err = os.Create(*outfilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer outfile.Close()

	w := bufio.NewWriter(outfile)
	defer w.Flush()

	fmt.Fprintln(w, commandsToString(commandsNoDupes))
}
