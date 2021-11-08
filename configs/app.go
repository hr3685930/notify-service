package configs

type App struct {
	Name      string `default:"notify"` //应用名
	Env       string `default:"local"`  //环境
	Debug     bool   `default:"true"`   //开启debug
	ErrReport string `default:"http://www.baidu.com"`
}
