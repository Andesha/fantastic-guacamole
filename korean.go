package main

import (
    "time"
    "bufio"
    "fmt"
    "os"
    "math/rand"
    "encoding/json"
    "io/ioutil"
    "reflect"
    "strings"
)

type Vocab struct {
    Days []Word `json:"days"`
    Sino []Word `json:"sino"`
    SinoTens []Word `json:"sino-tens"`
    Native []Word `json:"native"`
}

type Word struct {
    Korean string `json:"korean"`
    English string `json:"english"`
}

func (word Word) Quiz() bool {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("What is '%v' in english?: ", word.Korean)
    answer, _ := reader.ReadString('\n')
    answer = answer[:len(answer)-1] // Trim whitespace
    if answer == word.English {
        fmt.Println("Correct!")
        return true
    } else {
        fmt.Println("Wrong... Should be: ", word.English)
        return false
    }
}

func QuizLoop(words []Word){
    correct := 0
    totalQuestions := len(words)
    for _, value := range rand.Perm(totalQuestions) {
		if words[value].Quiz() {
            correct++
        }
	}
    fmt.Printf("Correct answers: %d\nTotal Questions: %d\n", correct, totalQuestions)
}

func main() {
    // Load file
    jsonFile, err := os.Open("vocab.json")
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    // Build structs
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var vocab Vocab
    json.Unmarshal(byteValue, &vocab)

    // Random seed set
    rand.Seed(time.Now().UnixNano())

    // Vocab field string builder
    s := reflect.ValueOf(&vocab).Elem()
    typeOfT := s.Type()
    var sb strings.Builder
    for i := 0; i < s.NumField(); i++ {
        if i != 0 {
            sb.WriteString(",")
        }
        sb.WriteString(strings.ToLower(typeOfT.Field(i).Name))
    }
    sb.WriteString(",scramble")

    // Menu system
    reader := bufio.NewReader(os.Stdin)
    quiz := true
    for quiz {
        fmt.Println("Select a quiz type: ", sb.String())
        answer, _ := reader.ReadString('\n')
        answer = answer[:len(answer)-1] // Trim whitespace
        switch answer {
        case "sino":
            QuizLoop(vocab.Sino)
        case "sinotens":
            QuizLoop(vocab.SinoTens)
        case "native":
            QuizLoop(vocab.Native)
        case "days":
            QuizLoop(vocab.Days)
        case "scramble":
            // v := reflect.ValueOf(vocab)
            // var collect []Word
            // values := make([]interface{}, v.NumField())
            // for i := 0; i < v.NumField(); i++ {
            //     values[i] = v.Field(i).Interface()
            //     for j:=0; j < len(values[i]); j++ {
            //         fmt.Println(values[i][j])
            //         // collect = append(collect,)
            //     }
            // }
            // fmt.Println(values)
            collect := append(vocab.Days, vocab.Sino...)
            collect = append(collect, vocab.SinoTens...)
            collect = append(collect, vocab.Native...)
            QuizLoop(collect)
        default:
            fmt.Println("Exiting...")
            quiz = false
        }
    }
}
