package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Book struct {
	ID                     int
	Title, Author          string
	Year, Available, Total int
}
type Reader struct {
	ID   int
	Name string
}
type Loan struct {
	ID, BookID, ReaderID int
	Active               bool
}

var books []Book
var readers []Reader
var loans []Loan
var bookIndex = map[int]int{}
var activeLoans = map[int][]int{}
var days = [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}

func reportHeader() {
	fmt.Println("Недельный отчёт:", strings.Join(days[:], ", "))
}

func findBook(id int) (*Book, error) {
	index, success := bookIndex[id]
	if !success {
		return nil, errors.New("книга не найдена")
	}
	return &books[index], nil
}

func findReader(id int) (*Reader, error) {
	for index := range readers {
		if readers[index].ID == id {
			return &readers[index], nil
		}
	}
	return nil, errors.New("читатель не найден")
}

// Уменьшаем остаток и сохраняем выдачу.
func borrowBook(book *Book, readerID int) error {
	if book.Available <= 0 {
		return errors.New("нет доступных экземпляров")
	}
	book.Available = book.Available - 1
	id := len(loans) + 1
	loans = append(loans, Loan{id, book.ID, readerID, true})
	activeLoans[readerID] = append(activeLoans[readerID], id)
	return nil
}

// Возвращаем книгу только при наличии активной выдачи.
func returnBook(book *Book, readerID int) error {
	for index, id := range activeLoans[readerID] {
		loan := &loans[id-1]
		if loan.BookID == book.ID && loan.Active {
			loan.Active = false
			book.Available = book.Available + 1
			ids := activeLoans[readerID]
			activeLoans[readerID] = append(ids[:index], ids[index+1:]...)
			return nil
		}
	}
	return errors.New("читатель не получал эту книгу или уже вернул её")
}

func searchBooks(query string) []Book {
	var found []Book
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), query) || strings.Contains(strings.ToLower(book.Author), query) {
			found = append(found, book)
		}
	}
	return found
}

func statistics() {
	reportHeader()
	available, total := 0, 0
	counts := map[int]int{}
	for _, book := range books {
		available = available + book.Available
		total = total + book.Total
	}
	for _, loan := range loans {
		counts[loan.BookID] = counts[loan.BookID] + 1
	}
	fmt.Printf("Наименований: %d\nВсего экземпляров: %d\nДоступно: %d\nВыдано: %d\n", len(books), total, available, total-available)
	max := 0
	var popular *Book
	for index := range books {
		if counts[books[index].ID] > max {
			popular = &books[index]
			max = counts[books[index].ID]
		}
	}
	if popular != nil {
		fmt.Printf("Самая популярная: %s, выдач: %d\n", popular.Title, max)
	}
}

var input = bufio.NewScanner(os.Stdin)

func read(prompt string) string {
	fmt.Print(prompt)
	if !input.Scan() {
		return "0"
	}
	return strings.TrimSpace(input.Text())
}

func number(prompt string) (int, error) {
	return strconv.Atoi(read(prompt))
}

func main() {
	for {
		fmt.Println("")
		fmt.Println("1. Добавить книгу")
		fmt.Println("2. Зарегистрировать читателя")
		fmt.Println("3. Выдать книгу")
		fmt.Println("4. Вернуть книгу")
		fmt.Println("5. Найти книгу")
		fmt.Println("6. Книги читателя")
		fmt.Println("7. Статистика")
		fmt.Println("0. Выход")
		choice := read("Действие: ")
		switch choice {
		case "1":
			id, operationError := number("ID книги: ")
			if operationError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			if _, success := bookIndex[id]; success {
				fmt.Println("ID уже существует")
				continue
			}
			title := read("Название: ")
			author := read("Автор: ")
			year, operationError := number("Год: ")
			if operationError != nil {
				fmt.Println("Неверный год")
				continue
			}
			count, operationError := number("Количество: ")
			if operationError != nil || count < 0 || title == "" || author == "" {
				fmt.Println("Неверные данные")
				continue
			}
			bookIndex[id] = len(books)
			books = append(books, Book{id, title, author, year, count, count})
			fmt.Println("Книга добавлена")
		case "2":
			id, operationError := number("ID читателя: ")
			if operationError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			if _, operationError = findReader(id); operationError == nil {
				fmt.Println("ID уже существует")
				continue
			}
			name := read("Имя: ")
			if name == "" {
				fmt.Println("Имя обязательно")
				continue
			}
			readers = append(readers, Reader{id, name})
			fmt.Println("Читатель добавлен")
		case "3", "4":
			readerID, operationError := number("ID читателя: ")
			if operationError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			reader, operationError := findReader(readerID)
			if operationError != nil {
				fmt.Println(operationError)
				continue
			}
			bookID, operationError := number("ID книги: ")
			if operationError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			book, operationError := findBook(bookID)
			if operationError != nil {
				fmt.Println(operationError)
				continue
			}
			if choice == "3" {
				operationError = borrowBook(book, readerID)
			} else {
				operationError = returnBook(book, readerID)
			}
			if operationError != nil {
				fmt.Println(operationError)
			} else {
				fmt.Printf("Успешно: %s, %s, доступно: %d\n", book.Title, reader.Name, book.Available)
			}
		case "5":
			query := strings.ToLower(read("Название или автор: "))
			found := searchBooks(query)
			if len(found) == 0 {
				fmt.Println("Книги не найдены")
			}
			for _, book := range found {
				fmt.Printf("#%d %s — %s (%d), доступно %d\n", book.ID, book.Title, book.Author, book.Year, book.Available)
			}
		case "6":
			readerID, operationError := number("ID читателя: ")
			if operationError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			if _, operationError = findReader(readerID); operationError != nil {
				fmt.Println(operationError)
				continue
			}
			if len(activeLoans[readerID]) == 0 {
				fmt.Println("Выданных книг нет")
			}
			for _, id := range activeLoans[readerID] {
				book, _ := findBook(loans[id-1].BookID)
				fmt.Printf("#%d %s\n", book.ID, book.Title)
			}
		case "7":
			statistics()
		case "0":
			return
		default:
			fmt.Println("Неизвестный пункт")
		}
	}
}
