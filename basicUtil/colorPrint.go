package basicUtil

import (
	"fmt"
	//"log"

)
//var PrintMessage = true

//func PrintlnBlue(a ...interface{}) {
//	//if PrintMessage {
//		fmt.Print(ANSI_COLOR_LIGHT_BLUE)
//		fmt.Print(a...)
//		fmt.Println(ANSI_COLOR_LIGHT_RESET)
//	//}
//	//log.Output(2, fmt.Sprintln(a...))
//}

func PrintlnRed(a ...interface{}) {
	fmt.Print(ANSI_COLOR_LIGHT_RED)
	fmt.Print(a...)
	fmt.Println(ANSI_COLOR_LIGHT_RESET)
}

func PrintlnYellow(a ...interface{}) {
	fmt.Print(ANSI_COLOR_LIGHT_YELLOW)
	fmt.Print(a...)
	fmt.Println(ANSI_COLOR_LIGHT_RESET)
}

func PrintlnMegenta(a ...interface{}) {
	fmt.Print(ANSI_COLOR_LIGHT_MAGENTA)
	fmt.Print(a...)
	fmt.Println(ANSI_COLOR_LIGHT_RESET)
}

func PrintlnGreen(a ...interface{}) {
	//if PrintMessage {
		fmt.Print(ANSI_COLOR_LIGHT_GREEN)
		fmt.Print(a...)
		fmt.Println(ANSI_COLOR_LIGHT_RESET)
	//}
	//log.Output(2, fmt.Sprintln(a...))
}

func PrintlnCyan(a ...interface{}) {
	//if PrintMessage {
		fmt.Print(ANSI_COLOR_LIGHT_CYAN)
		fmt.Print(a...)
		fmt.Println(ANSI_COLOR_LIGHT_RESET)
	//}
	//log.Output(2, fmt.Sprintln(a...))
}


