package converters

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	"io/ioutil"
	"unicode"

	"github.com/zxjsdp/bioinfo-go/utils"
)

const freeSpaceNum int = 4

var (
	fileInfo *os.FileInfo
	err      error
)

func ExtractSpeciesFromFastaFile(fastaFile string) []Species {
	var species []Species
	var title, sequence, line string
	file, err := os.Open(fastaFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line != "" {
			if strings.HasPrefix(line, ">") {
				if sequence != "" {
					species = append(species, Species{title, sequence})
					sequence = ""
				}
				title = strings.TrimLeft(line, ">")
				title = replaceBlankChars(title)
			} else {
				sequence = sequence + strings.TrimSpace(line)
			}
		}
	}
	if title != "" && sequence != "" {
		species = append(species, Species{title, sequence})
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return species
}

func GeneratePhylip(species []Species, outputFile string) {
	if species == nil {
		log.Panic("No species!")
		return
	}
	var phylipLines []string
	longestNameLength := getLongestNameLength(species)
	speciesNum := len(species)
	charNum := len(species[0].sequence)

	// Add (speciesNum  charNum) to top line of output Phylip file
	phylipLines = append(phylipLines, fmt.Sprintf("%d  %d", speciesNum, charNum))

	// Add (species  sequence) to each line
	for _, each := range species {
		spacesForCurrentSpecies := generateSpacesForAlignment(longestNameLength, len(each.name))
		phylipLines = append(phylipLines, each.name + spacesForCurrentSpecies + each.sequence)
	}
	phylipContent := strings.Join(phylipLines, "\n")

	err := ioutil.WriteFile(outputFile, []byte(phylipContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func replaceBlankChars(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, str)
}

func getLongestNameLength(species []Species) int {
	longestNameLength := 0
	for _, each := range species {
		if len(each.name) > longestNameLength {
			longestNameLength = len(each.name)
		}
	}
	return longestNameLength
}

func generateSpacesForAlignment(longestNameLength, currentNameLen int) string {
	spaceNum := longestNameLength - currentNameLen + freeSpaceNum
	return utils.GenerateSpaces(spaceNum)
}

