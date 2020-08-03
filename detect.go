// This file will receive a Pokemon Image and will return the appropriate pokemon.
// This will get called from jokercord.go

package main

import (
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/corona10/goimagehash"

	"gopkg.in/yaml.v2"
)

// BEGIN structs
// ENDOF structs

// BEGIN function definition
// receive grabs url of Pokemon picture
func receive(url string) string {
	pokemons := make(map[string]string)
	err := Download("images/template.jpg", url)
	logErr(err)
	hash := Hash("images/template.jpg")
	readPokemonList(pokemons)
	pokemonName := Compare(hash, pokemons)
	return (pokemonName)
}

// Download grabs Pokemon Picture from recieve url
func Download(path string, url string) error {
	response, err := http.Get(url)
	logErr(err)
	var output *os.File

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logErr(err)
		output, err = os.Create(path)
		logErr(err)
	} else {
		err = os.Remove(path)
		logErr(err)
		output, err = os.Create(path)
		logErr(err)
	}
	// copy contents of response to output
	_, err = io.Copy(output, response.Body)
	defer response.Body.Close()
	defer output.Close()
	return err
}

// readPokemonList reads hash list
func readPokemonList(pokemonStruct map[string]string) {
	reader, err := ioutil.ReadFile("config/hashes.yaml")
	logErr(err)
	yaml.Unmarshal(reader, pokemonStruct)
}

// Compare checks hash to hash list
func Compare(hash string, pokemonStruct map[string]string) string {
	var name string
	hash = strings.Replace(hash, "p:", "", 1)
	for pokemon, pokemonHash := range pokemonStruct {
		if pokemonHash == hash {
			name = pokemon
		}
	}
	return name

}

// Hash grabs value from Download
func Hash(path string) string {
	output, err := os.Open(path)
	logErr(err)
	imageFile, _, err := image.Decode(output)
	if err != nil {
		log.Panic("Couldn't read image file")
	}
	hash, err := goimagehash.PerceptionHash(imageFile)
	if err != nil {
		log.Panic("Couldn't get hash")
	}
	return hash.ToString()
}

// ENDOF function definition
