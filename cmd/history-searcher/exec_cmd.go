package main

import (
	"bufio"
	"os"
	"text/template"
)

// execCmd executes the command in the current process
// func execCmd(cmdStr string) {
// 	args := strings.Split(cmdStr, " ")
// 	env := os.Environ()

// 	binary, lookErr := exec.LookPath(args[0])
// 	if lookErr != nil {
// 		panic(lookErr)
// 	}

// 	execErr := syscall.Exec(binary, args, env)
// 	if execErr != nil {
// 		panic(execErr)
// 	}
// }

func writeCmdToFile(cmdStr, outPutFile string) {
	temp = template.Must(template.ParseFS(myTemplates, "templates/cmd-exe.tmpl"))

	wr, err := os.OpenFile(outPutFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	defer wr.Close()
	w := bufio.NewWriter(wr)
	err = temp.Execute(w, cmdStr)

	if err != nil {
		panic(err)
	}
	w.Flush()
}

/*
Mode options?
- Just print out the command, the user can copy and paste it
- Execute the command if it is not a bash built-in
- Write the command to a file, allowing the user to execute it later with source
- Create a wrapper function that executes the searcher and then runs source on the command
*/
