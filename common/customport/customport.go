package customport

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	// Ports map[string][]int = map[string][]int{
	// 	"default": {80, 443, 8080, 8081, 8083},
	// 	"custom":  {},
	// }
	DefaultPorts = []int{80, 443, 8080, 8081, 8083}
)

type Ports []int

func (c *Ports) String() string {
	return "Custom Ports"
}

func (c *Ports) Set(value string) error {
	*c = []int{}
	potentialPorts := strings.Split(strings.Trim(value, ","), ",")
	for _, potentialPort := range potentialPorts {
		potentialRange := strings.Split(potentialPort, "-")
		if len(potentialRange) < 2 {
			if p, err := strconv.Atoi(potentialPort); err == nil {
				*c = append(*c, p)
			} else {
				fmt.Printf("Could not cast port to integer, your value: %s, resulting error %s. Skipping it\n",
					potentialPort, err.Error())
			}
		} else {
			var lowP, highP int
			lowP, err := strconv.Atoi(strings.Trim(potentialRange[0], ", "))
			if err != nil {
				fmt.Printf("Could not cast first port of your port range(%s) to integer, your value: %s, resulting error %s. Skipping it\n",
					potentialPort, potentialRange[0], err.Error())
				continue
			}
			highP, err = strconv.Atoi(strings.Trim(potentialRange[1], ", "))
			if err != nil {
				fmt.Printf("Could not cast last port of your port range(%s) to integer, "+
					"your value: %s, resulting error %s. Skipping it\n",
					potentialPort, potentialRange[1], err.Error())
				continue
			}

			if lowP > highP {
				fmt.Printf("first value of port range should be lower than the last part port "+
					"in that range, your range: [%d, %d]. Skipping it\n",
					lowP, highP)
				continue
			}

			for i := lowP; i <= highP; i++ {
				*c = append(*c, i)
			}
		}
	}
	// fmt.Println
	return nil
}
