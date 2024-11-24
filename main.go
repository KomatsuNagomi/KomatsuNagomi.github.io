package main

import (
    "fmt"
    "github.com/labstack/echo/v4"
    "net/http"
    "os"
    "strconv"
    "errors"
    "strings"
)

type Data struct {
    Num string `json:"num" form:"num"`
}

func showForm(c echo.Context) error {
    htmlContent, err := os.ReadFile("index.html")
    if err != nil {
        return c.String(http.StatusInternalServerError, "unable to read the HTML file") 
    }
    return c.HTML(http.StatusOK, string(htmlContent))
}

func calculator(c echo.Context) error {
    var input Data // 1 + 2 * 3

    if err := c.Bind(&input); err != nil {
        return c.String(http.StatusBadRequest, "Bad request")
    }

    ans, err := parse(input.Num)

    if err == nil {
        return c.HTML(http.StatusOK, fmt.Sprintf("ans : %d", ans))
    } else {
        errMessage := err.Error()
        parts := strings.Split(errMessage, ": ")
        return c.String(http.StatusBadRequest, parts[0]) 
    }
}

func parse(input string) (int, error) {
    input = strings.ReplaceAll(input, " ", "")

    if !isJadgeEndNum(input) {
        return 0, errors.New("the last character is not a number")  
    }

    nums, operators, err := getNumsAndOperators(input)
    if err != nil {
        return 0, fmt.Errorf("%w",err)
    }

    if len(nums) != len(operators)+1 {
        return 0, errors.New("the number of operators and numbers is mismatched") 
    }

    var newNums []int
    var newOperators []string
    currentNum := nums[0]

    for i, operator := range operators {
        rightNum := nums[i+1]
        if operator == "*" {
            currentNum *= rightNum
        } else if operator == "/" {
            if rightNum == 0 {
                return 0, errors.New("cannot divide by zero")  
            }
            currentNum /= rightNum
        } else {
            newNums = append(newNums, currentNum)
            newOperators = append(newOperators, operator)
            currentNum = rightNum
        }
    }

    newNums = append(newNums, currentNum)

    for i := 0; i < len(newOperators); i++ {
        leftNum := newNums[i]
        operator := newOperators[i]
        rightNum := newNums[i+1]

        var result int
        if operator == "+" {
            result = leftNum + rightNum
        } else if operator == "-" {
            result = leftNum - rightNum
        }

        newNums[i+1] = result
    }

    return newNums[len(newNums)-1], nil
}

func isJadgeEndNum(input string) bool {
    var firstString = string(input[0])
    var endString = string(input[len(input)-1])

    _, firstErr := strconv.Atoi(firstString)
    _, endErr := strconv.Atoi(endString)

    return firstErr == nil && endErr == nil
}

func getNumsAndOperators(input string) ([]int, []string, error) {
    var temp string
    var nums []int
    var operators []string

    OPERATORS := map[string]bool{
        "+": true,
        "-": true,
        "*": true,
        "/": true,
    }

    for _, str := range input {
        char := string(str)
        _, err := strconv.Atoi(char)

        if err == nil {
            temp = temp + char
            continue
        } 

        if temp != "" && err != nil && OPERATORS[char] {
            num, _ := strconv.Atoi(temp)
            operators = append(operators, char)
            nums = append(nums, num)
            temp = ""
            continue
        }
        
        return nil, nil, errors.New("invalid character in the input") 
    }
    
    if temp != "" {
        num, err := strconv.Atoi(temp)
        if err != nil {
            return nil, nil, errors.New("the expression ends with an operator") 
        }
        nums = append(nums, num)
    }
    return nums, operators, nil
}


func main() {
    e := echo.New()
    e.GET("/", showForm)

	e.POST("/add", calculator)

	e.Logger.Fatal(e.Start(":1323"))
}
