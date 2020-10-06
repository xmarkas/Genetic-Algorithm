package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

/**
* Population size that will static
 */
var populationSize int = 100

/**
* The top % of the fittest that will mate and create offspring.
 */
const matingPopulation = 20

/**
* The % of adults from the current generation that will
* survive to the next generation
 */
var survivingAdults int = 10

/**
* The % of the genes that will have a probability of being
* mutated.
 */
var mutationPercent float32 = .10

var population = make([]Individual, populationSize)
var target string
var genes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*() -+=?':;<>.{}[]/"

func main() {
	rand.Seed(time.Now().UnixNano())
	var numOfSamples = 1

	// User interface
	flag := true
	for flag {
		target = targetInput()
		populationSize = populationInput()
		survivingAdults = survivingInput()

		// Confirm choices
		fmt.Println("\n-------------------------")
		fmt.Println("- Target string: ", target)
		fmt.Println("- Population size: ", populationSize)
		fmt.Println("- % of surviving adults: ", survivingAdults)
		fmt.Println("-------------------------")
		fmt.Println("is this correct (y/n)?")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Text() == "y" {
			flag = false
		}
	}

	var test Sample
	for i := 0; i < numOfSamples; i++ {
		test.generations = append(test.generations, runSample())
	}

	fmt.Println(test)

	var totalGen int = 0
	for n := 0; n < len(test.generations); n++ {
		totalGen += test.generations[n]
	}
	fmt.Println("Average generations till success: ", totalGen/numOfSamples)

}

type UserInputs struct {
	target string "Enter a string to target: "
}

func targetInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter a string to target: ")
	scanner.Scan()
	return scanner.Text()
}

func populationInput() int {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\nPopulation size (default 100): ")

	scanner.Scan()
	size, _ := strconv.Atoi(scanner.Text())
	for len(scanner.Text()) != 0 && size == 0 {
		fmt.Println("Invalid value. Must be a whole number")
		scanner.Scan()
		size, _ = strconv.Atoi(scanner.Text())
	}

	if len(scanner.Text()) == 0 {
		return 100
	} else {
		return size
	}

}

func survivingInput() int {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\nPercent of fittest adults that will go to\nthe next generation (default 10%)\nEnter as a whole number: ")

	scanner.Scan()
	size, _ := strconv.Atoi(scanner.Text())
	for len(scanner.Text()) != 0 && size == 0 {
		fmt.Println("Invalid value. Must be a whole number")
		scanner.Scan()
		size, _ = strconv.Atoi(scanner.Text())
	}

	if len(scanner.Text()) == 0 {
		return 10
	} else {
		return size
	}
}

func runSample() int {
	var found bool = false
	var generation int = 1

	// Create 1st generation
	newGeneration(0)

	for !found {
		// Sort population from highest fitness to lowest
		sort.Slice(population[:], func(i, j int) bool {
			return population[i].fitnessScore < population[j].fitnessScore
		})

		// If target exists in population
		if population[0].fitnessScore == 0 {
			fmt.Println("Generation: ", generation, "Individual: ", population[0])
			found = true
			break
		}

		// % of fittest population will mate to produce offspring
		for k := 0; k < (populationSize - (populationSize * survivingAdults / 100)); k++ {
			var parent1 Individual = population[rand.Intn(populationSize*matingPopulation/100)]
			var parent2 Individual = population[rand.Intn(populationSize*matingPopulation/100)]

			var offspring Individual = parent1.mateIndividual(parent2)

			// % of existing population goes to next generation
			population[k+(populationSize*survivingAdults/100)] = offspring
		}

		fmt.Println("Generation: ", generation, "Fitness: ", population[0].fitnessScore, population[0].genome)

		generation++
	}

	return generation
}

type Individual struct {
	genome       string
	fitnessScore int
}

func (ind *Individual) mateIndividual(mate Individual) Individual {
	var offspringChromosone string = ""

	for i := 0; i < len(target); i++ {
		// random probability
		var prob float32 = rand.Float32()

		// if prob is less than 45% insert allele from parent1
		if prob < ((1 - mutationPercent) / 2) {
			offspringChromosone += string(ind.genome[i])
		} else if prob < (1 - mutationPercent) { // if prob is between 45% and 90% insert allele from parent2
			offspringChromosone += string(mate.genome[i])
		} else { // Otherwise random gene for diversity
			offspringChromosone += string(mutatedGene())
		}
	}

	var offspring Individual
	offspring.genome = offspringChromosone
	offspring.fitnessScore = getFitnessScore(offspring.genome)

	return offspring
}

type Sample struct {
	sampleSize  int
	generations []int
}

func newGeneration(begin int) {
	population = make([]Individual, populationSize)
	for i := begin; i < populationSize; i++ {
		var newIndividual Individual
		newIndividual.genome = createChromosome()
		newIndividual.fitnessScore = getFitnessScore(newIndividual.genome)
		population[i] = newIndividual
	}
}

func mutatedGene() byte {
	return genes[rand.Intn(len(genes))]
}

func createChromosome() string {
	var gene string = ""

	for i := 0; i < len(target); i++ {
		gene += string(mutatedGene())
	}

	return gene
}

func getFitnessScore(g string) int {
	var score int = 0
	for i := 0; i < len(target); i++ {
		if g[i] != target[i] {
			score++
		}
	}

	return score
}
