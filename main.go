package main

import (
  "bufio"
  "bytes"
  "errors"
  "fmt"
  "os/exec"
  "regexp"
  "strings"

  //"fmt"
  "log"
  "os"

  "github.com/urfave/cli/v2"
)

func runCommand(command string) string {
  fmt.Println(command)

  args := strings.Split(command, " ")
  cmd := exec.Command(args[0], args[1:]...)

  currentDirectory,_ := os.Getwd()
  cmd.Dir = currentDirectory

  var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
  if err != nil {
    fmt.Println(err.Error())
    return err.Error()
  }

  //fmt.Println(out.String())
  return out.String()
}


func pythonExists(pythonObject PythonObject) bool {
	outputs := runCommand(pythonObject.pythonCommand + " --version")
	if strings.Contains(outputs, "not found") {
      return false
    } else {
      return true
    }
}

func extractNumbersFromString(text string) []string {
  re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

  return re.FindAllString(text, -1)
}

type PythonObject struct {
  valid bool
  versionCode string
  pythonCommand string
  pipCommand string
}

func parseThePythonTag(theFirstLineOfCode string) PythonObject {
    pythonObject := PythonObject{
      valid: false,
      versionCode: "",
      pythonCommand: "",
      pipCommand: "",
    }

	if theFirstLineOfCode == "" {
	  return pythonObject
    }

    if len(theFirstLineOfCode) < len("# python") {
      return pythonObject
    }

    if theFirstLineOfCode[0:7] != "#python" && theFirstLineOfCode[0:8] != "# python" {
      return pythonObject
    }

    listOfFloatNumbers := extractNumbersFromString(theFirstLineOfCode)
    if len(listOfFloatNumbers) < 1 {
      return pythonObject
    } else {
      var pythonVersionCode string = listOfFloatNumbers[0]
      pythonObject.versionCode = pythonVersionCode
      pythonObject.pythonCommand = "python" + pythonVersionCode
      pythonObject.pipCommand = "pip" + pythonVersionCode
      return pythonObject
    }
}

func readTheFirstLineOfCodeFromAFile(textFilePath string) string {
  f, err := os.Open(textFilePath)
  if err != nil {
    log.Fatalln(err)
  }
  defer f.Close()

  scanner := bufio.NewScanner(f)
  var line int = 0
  for scanner.Scan() {
    if line == 0 {
      theFirstLine := scanner.Text()
      return theFirstLine
    }
    line++
  }

  return ""
}


func takeThePythonFilePathAndDoSomething(pythonFilePath string) error {
  if pythonFilePath == "" {
    return errors.New("no python file given")
  }

  // check if that file is python file
  if len(pythonFilePath) < 3 {
    return errors.New("not a python file")
  }
  if pythonFilePath[len(pythonFilePath)-3:] != ".py" {
    return errors.New("not a python file")
  }

  // check if the python file exist
  _, err := os.Stat(pythonFilePath)
  if err != nil {
    return err
  }

  theFirstLine := readTheFirstLineOfCodeFromAFile(pythonFilePath)
  pythonObject := parseThePythonTag(theFirstLine)

  if pythonExists(pythonObject) == true {
  	// python run
    outputs := runCommand(pythonObject.pythonCommand + " " + pythonFilePath)
    fmt.Println(outputs)
    return nil
  } else {
    // install that version of python
  	command := fmt.Sprintf("sudo apt install %s", pythonObject.pythonCommand)
    //command := fmt.Sprintf("brew install %s", pythonObject.pythonCommand)
    outputs := runCommand(command)
    fmt.Println(outputs)

    outputs = runCommand(pythonObject.pythonCommand + " " + pythonFilePath)
    fmt.Println(outputs)
  }

  return nil
}

func main() {
  app := &cli.App{
    Action: func(c *cli.Context) error {
      var python_file_path string = c.Args().Get(0)
      err := takeThePythonFilePathAndDoSomething(python_file_path)
      //fmt.Printf("Hello %v\n\n", )
      return err
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
