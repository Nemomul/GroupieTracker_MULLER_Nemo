package main

//Import les fichiers nécessaires
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Structure pour stocker les informations d'un personnage
type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Origin  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Image   string   `json:"image"`
	Episode []string `json:"episode"`
	URL     string   `json:"url"`
}

// Fonction principale
func main() {

	// Configuration des fichiers statiques
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	// Route pour la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			id = "1"
		}
		character, err := getCharacter(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Exécution du template pour afficher le personnage
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, character)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Route pour la page de test
	http.HandleFunc("/seconde", func(w http.ResponseWriter, r *http.Request) {
		characters, err := getAllCharacters()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Exécution du template pour afficher la liste des personnages
		tmpl, err := template.ParseFiles("static/html/seconde.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, characters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	//Route pour les cartes des personnages
	http.HandleFunc("/Perso", func(w http.ResponseWriter, r *http.Request) {
		idPerso := strings.TrimPrefix(r.URL.RequestURI(), "/Perso?id=")
		p1, _ := getCharacter(idPerso)

		tmpl, _ := template.ParseFiles("static/html/cards.html")
		tmpl.Execute(w, p1)

	})

	// Démarrage du serveur web
	http.ListenAndServe(":80", nil)
}

// Récupération d'un personnage par son identifiant
func getCharacter(id string) (*Character, error) {

	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var character Character
	err = json.NewDecoder(resp.Body).Decode(&character)
	if err != nil {
		return nil, err
	}

	return &character, nil
}
func getAllCharacters() ([]Character, error) {
	url := "https://rickandmortyapi.com/api/character"
	var characters []Character

	for {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var pageCharacters struct {
			Info struct {
				Next string `json:"next"`
			} `json:"info"`
			Results []Character `json:"results"`
		}

		err = json.NewDecoder(resp.Body).Decode(&pageCharacters)
		if err != nil {
			return nil, err
		}

		characters = append(characters, pageCharacters.Results...)

		if pageCharacters.Info.Next == "" {
			break
		}

		url = pageCharacters.Info.Next
	}

	return characters, nil
}
