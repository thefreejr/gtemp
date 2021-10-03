package main


import (
          "fmt"
      //  "flag"
        "os"
        "log"
        "time"
        "strings"
        "strconv"
        "regexp"
        "github.com/tkanos/gonfig"
        "github.com/stianeikeland/go-rpio/v4"
)

type pc struct {
	Temp_Limit float64
	GPIO_port string
  GPIO int
  Log_Path string
  Hysteresys float64
  Current_CPU_temp float64
  Fan_Enable bool
  CPU_Temp_Path string



}
var node pc = pc{}
var cfgFile string = "gtemp.conf"
err := rpio.Open()
var pinStat bool = true



func check(err error) {
  if err != nil {
    log.Fatal(err.Error())
  }
}

func readCfg(cfg string) (n pc) {

  err := gonfig.GetConf(cfg, &n)
  check(err)
  //рассчитываем номер порта GPIO
  n.GPIO,_ = convertGpioPort(n.GPIO_port)
  n.Fan_Enable = pinStat //         rpio.Pin(n).Read()
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


func logConfigure() {
  if node.Log_Path != ""{
    file, err := os.OpenFile(node.Log_Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    check(err)
    log.SetOutput(file)
  }

}

func doWork() {
  for {
    node.Current_CPU_temp = getCPUTemp(node.CPU_Temp_Path)
    if node.Current_CPU_temp >= node.Temp_Limit {
      fanControll(node.GPIO, true)

  } else if node.Current_CPU_temp <= node.Temp_Limit - node.Hysteresys {
      fanControll(node.GPIO, false)
    }
    time.Sleep(1 * time.Second)
  }
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
  func  fanControll(n int , command bool) {
    node.Fan_Enable = pinStat //         rpio.Pin(n).Read()
    if command && node.Fan_Enable {
      //rpio.Pin(n).High()
      log.Printf("Current CPU Temperature: %4.2f(%4.2f), fan is ON",node.Current_CPU_temp, node.Temp_Limit)
    }
    if !command && !node.Fan_Enable {
      //rpio.Pin(n).Low()
      log.Printf("Current CPU Temperature: %4.2f(%4.2f), fan is OFF",node.Current_CPU_temp, node.Temp_Limit)
    }
  }



func main() {
  //определяем модель и версию платы
//  getHW()

  // читаем конфиг из файла или из ключей
  node = readCfg(cfgFile)

  // настраиваем логирование
  logConfigure()

  // считываем температуру в цикле и управляем
  doWork()



  fmt.Println(node)

}
