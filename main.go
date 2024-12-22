package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Calc вычисляет арифметическое выражение, переданное в виде строки,
// и возвращает результат в виде float64 или ошибку, если выражение некорректно.
func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	return evaluatePostfix(postfix)
}

func tokenize(expr string) []string {
	var tokens []string
	var currentToken strings.Builder
	for _, char := range expr {
		switch char {
		case ' ':
			continue
		case '+', '-', '*', '/', '(', ')':
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		default:
			currentToken.WriteRune(char)
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}
	return tokens
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var operators []string

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, errors.New("Недостающие скобки")
			}
			// Удаляем "(" из стека
			operators = operators[:len(operators)-1]
		} else if isOperator(token) {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else {
			return nil, fmt.Errorf("Некорректный ввод")
		}
	}
	// Выгружаем оставшиеся операторы
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, errors.New("Недостающие скобки")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("Некорректный ввод")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("Деление на ноль")
				}
				stack = append(stack, a/b)
			default:
				return 0, fmt.Errorf("Неизвестный оператор: %s", token)
			}
		} else {
			return 0, fmt.Errorf("Некорректный ввод: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("Некорректный ввод")
	}
	return stack[0], nil
}

func isNumber(token string) bool {
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return true
	}
	return false
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// Для удобства работы с JSON определяем структуру входных данных
type requestData struct {
	Expression string `json:"expression"`
}

// И структуру для корректного ответа
type responseData struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// calculateHandler обрабатывает POST-запрос на /api/v1/calculate
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Принимаем только POST-запросы
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON
	var req requestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Ошибка декодирования → 422 Unprocessable Entity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(responseData{Error: "Expression is not valid"})
		return
	}

	// Вычисляем результат
	result, err := Calc(req.Expression)
	if err != nil {
		// Если в тексте ошибки что-то явно указывает на неверный ввод — отдаем 422
		if strings.Contains(err.Error(), "Некорректный ввод") ||
			strings.Contains(err.Error(), "Недостающие скобки") ||
			strings.Contains(err.Error(), "Деление на ноль") {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(responseData{Error: "Expression is not valid"})
		} else {
			// Иначе — 500 Internal Server Error
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseData{Error: "Internal server error"})
		}
		return
	}

	// Если всё успешно
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData{Result: fmt.Sprintf("%g", result)})
}

func main() {
	// Регистрируем наш обработчик на /api/v1/calculate
	http.HandleFunc("/api/v1/calculate", calculateHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
