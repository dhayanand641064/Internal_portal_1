package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
)

type Profile struct {
	Name         string
	RollNumber   int
	Email        string
	GitHub       string
	Role         string
	IsAdmin      bool
	Projects     []string
	CodingLinks  []string
	Skills       []string
	Interests    []string
	Achievements []string
}

func createProfile(session *gocql.Session, profile Profile) error {
	query := "INSERT INTO dhaya5 (column1_name, column2_roll_number, column3_email, column4_github, column5_role, column6_is_admin, column10_projects, column11_coding_links, column7_skills, column8_interests, column9_achievements) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return session.Query(query, profile.Name, profile.RollNumber, profile.Email, profile.GitHub, profile.Role, profile.IsAdmin, profile.Projects, profile.CodingLinks, profile.Skills, profile.Interests, profile.Achievements).Exec()
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "my_app"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer session.Close()

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {

		name := r.URL.Query().Get("name")

		getProfileHandler(w, r, session, name)
	})

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProfileHandler(w http.ResponseWriter, r *http.Request, session *gocql.Session, name string) {
	profile, err := getProfile(session, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, "Failed to marshal profile to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getProfile(session *gocql.Session, name string) (Profile, error) {
	var profile Profile
	query := "SELECT * FROM dhaya5 WHERE column1_name = ?"
	iter := session.Query(query, name).Iter()

	var result map[string]interface{}
	result = make(map[string]interface{})

	if iter.MapScan(result) {
		profile.Name = result["column1_name"].(string)
		profile.RollNumber = result["column2_roll_number"].(int)
		profile.Email = result["column3_email"].(string)
		profile.GitHub = result["column4_github"].(string)
		profile.Role = result["column5_role"].(string)
		profile.IsAdmin = result["column6_is_admin"].(bool)
		profile.Projects = result["column10_projects"].([]string)
		profile.CodingLinks = result["column11_coding_links"].([]string)
		profile.Skills = result["column7_skills"].([]string)
		profile.Interests = result["column8_interests"].([]string)
		profile.Achievements = result["column9_achievements"].([]string)

		return profile, nil
	}

	if err := iter.Close(); err != nil {
		return profile, fmt.Errorf("Failed to retrieve profile: %v", err)
	}

	return profile, fmt.Errorf("Profile not found")
}
