package main

import ("fmt"
      //  "flag"
        "os"
        //"log"
        "strings"
        "regexp"
        "github.com/tkanos/gonfig"
      //  "github.com/stianeikeland/go-rpio/v4"
)

type node struct {
	Temp_Limit int
	GPIO_port string
  Log_level string
  Hysteresys int
  Current_CPU_temp int
  Fan_Enable bool

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
  n.GPIO_port,_ = convertGpioPort(n.GPIO_port)
  return n
}

func convertGpioPort(s string) (t string, err error) {
  re, _ := regexp.Compile(`^gpio\d_[a-z]\d`)
  matched := re.MatchString(s)

  if matched {
    var s1 []string = strings.Split(s, "")
    m := map[string]string {
      "a": "1",
      "b": "2",
      "c": "3",
    }
    t = s1[4] + m[s1[6]] + s1[7]

  } else  {
    t =  s
  }

  return t, err
  }




func main() {
  //определяем модель и версию платы
//  getHW()

  // читаем конфиг из файла или из ключей
  node := readCfg(cfgFile)

  fmt.Println(node)

  // считываем температуру в цикле
  //getCPUTemp()

  // управляем
  //fanControll()


}
