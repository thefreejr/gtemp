package main


import (
          "fmt"
      //  "flag"
        "os"
        //"log"
        "strings"
        "strconv"
        "regexp"
        "github.com/tkanos/gonfig"
      //  "github.com/stianeikeland/go-rpio/v4"
)

type node struct {
	Temp_Limit int
	GPIO_port string
  GPIO int
  Log_level string
  Hysteresys int
  Current_CPU_temp float64
  Fan_Enable bool
  CPU_Temp_Path string


}
var cfgFile string = "gtemp.conf"

func check(err error) {
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
}

func readCfg(cfg string) (n node) {

  err := gonfig.GetConf(cfg, &n)
  check(err)
  //рассчитываем номер порта GPIO
  n.GPIO,_ = convertGpioPort(n.GPIO_port)
  return n
}

func convertGpioPort(s string) (t int, err error) {
  re, err := regexp.Compile(`^gpio\d_[a-z]\d`)
  matched := re.MatchString(s)

  if matched {
    var s1 []string = strings.Split(s, "")
    m := map[string]string {
      "a": "0",
      "b": "1",
      "c": "2",
      "d": "3",
    }
    t1,_ := strconv.Atoi(s1[4])
    t2,_ := strconv.Atoi(m[s1[6]])
    t3,_ := strconv.Atoi(s1[7])
    t =  t1*32 + t2*8 + t3
  } else  {
    t,err =  strconv.Atoi(s)
  }
  return t, err
}

func getCPUTemp(path string) (t float64) {
  bytes, err := os.ReadFile(path);
	check(err)
	fileText := string(bytes[:]);
  re, err := regexp.Compile(`^\d+`)
  t,err = strconv.ParseFloat(re.FindString(fileText), 64)
  t = t/1000
  check(err)
  return t
}



func main() {
  //определяем модель и версию платы
//  getHW()

  // читаем конфиг из файла или из ключей
  node := readCfg(cfgFile)



  // считываем температуру в цикле
  node. Current_CPU_temp = getCPUTemp(node.CPU_Temp_Path)

  // управляем
  //fanControll()
  fmt.Println(node)

}
