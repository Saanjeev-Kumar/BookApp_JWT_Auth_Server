package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"userType"`
}

type Book struct {
	Name            string `json:"name"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}

var sampleUsers = []User{
	{Username: "admin", Password: "adminpass", UserType: "admin"},
	{Username: "user", Password: "userpass", UserType: "regular"},
}

var mySigningKey = []byte("secret")

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/home", Home).Methods("GET")
	r.HandleFunc("/addBook", AddBook).Methods("POST")
	r.HandleFunc("/deleteBook", DeleteBook).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	for _, u := range sampleUsers {
		if user.Username == u.Username && user.Password == u.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": user.Username,
				"userType": u.UserType,
			})

			tokenString, error := token.SignedString(mySigningKey)
			if error != nil {
				fmt.Println(error)
			}
			json.NewEncoder(w).Encode(tokenString)
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	userType := claims["userType"].(string)
	fmt.Println(userType)
	books := getBooksForUserType(userType)
	//books := "Hello world book"
	json.NewEncoder(w).Encode(books)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if claims["userType"] != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var book Book
	json.NewEncoder(w).Encode("New book added")
	_ = json.NewDecoder(r.Body).Decode(&book)

	if book.Name == "" || book.Author == "" || book.PublicationYear <= 0 {
		http.Error(w, "Invalid book data", http.StatusBadRequest)
		return
	}

	file, _ := os.OpenFile("regularUser.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{book.Name, book.Author, fmt.Sprintf("%d", book.PublicationYear)})
}

// func DeleteBook(w http.ResponseWriter, r *http.Request) {
// 	tokenString := r.Header.Get("Authorization")
// 	claims := jwt.MapClaims{}
// 	_, _ = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return mySigningKey, nil
// 	})

// 	if claims["userType"] != "admin" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	bookName := r.URL.Query().Get("name")
// 	if bookName == "" {
// 		http.Error(w, "Invalid book name", http.StatusBadRequest)
// 		return
// 	}

// 	books := getBooksForUserType("regular")
// 	for i, book := range books {
// 		if strings.ToLower(book.Name) == strings.ToLower(bookName) {
// 			books = append(books[:i], books[i+1:]...)
// 			break
// 		}
// 	}

// 	file, _ := os.Create("regularUser.csv")
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

//		for _, book := range books {
//			writer.Write([]string{book.Name, book.Author, fmt.Sprintf("%d", book.PublicationYear)})
//		}
//	}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if claims["userType"] != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("delete check auth done")
	bookName := r.URL.Query().Get("name")
	fmt.Println("bookname", bookName)
	if bookName == "" {
		http.Error(w, "Invalid book name", http.StatusBadRequest)
		return
	}

	books := getBooksForUserType("regular")
	var updatedBooks []Book
	for _, book := range books {
		if book.Name != bookName {
			updatedBooks = append(updatedBooks, book)
		}
	}

	file, _ := os.Create("regularUser.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, book := range updatedBooks {
		writer.Write([]string{book.Name, book.Author, fmt.Sprintf("%d", book.PublicationYear)})
	}
}

func getBooksForUserType(userType string) []Book {
	// var books []Book
	// fmt.Println("Inside function of homw")
	// file, _ := os.Open("regularUser.csv")
	// defer file.Close()
	// fmt.Println(file)
	// reader := csv.NewReader(file)
	// fmt.Println(reader)
	// for {
	// 	record, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if len(record) >= 2 {
	// 		books = append(books, Book{Name: record[0], Author: record[1]})
	// 	}
	// 	// books = append(books, Book{Name: record[0], Author: record[1]})
	// 	fmt.Println(books)
	// }

	// if userType == "admin" {
	// 	file, _ = os.Open("adminUser.csv")
	// 	defer file.Close()
	// 	reader = csv.NewReader(file)
	// 	for {
	// 		record, err := reader.Read()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if len(record) >= 2 {
	// 			books = append(books, Book{Name: record[0], Author: record[1]})
	// 		}
	// 		//books = append(books, Book{Name: record[0], Author: record[1]})
	// 	}
	// }

	// return books
	var books []Book

	file, _ := os.Open("regularUser.csv")
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading from CSV file:", err)
		return books
	}
	for _, record := range records {
		if len(record) >= 2 {
			books = append(books, Book{Name: record[0], Author: record[1]})
		}
	}

	if userType == "admin" {
		file, _ = os.Open("adminUser.csv")
		defer file.Close()
		reader = csv.NewReader(file)
		records, err = reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading from CSV file:", err)
			return books
		}
		for _, record := range records {
			if len(record) >= 2 {
				books = append(books, Book{Name: record[0], Author: record[1]})
			}
		}
	}

	return books
}
