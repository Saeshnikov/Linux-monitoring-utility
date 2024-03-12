package lsofLayer

func LsofExec() []string {
	cmd := exec.Command("/usr/bin/lsof")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	outScanner := bufio.NewScanner(stdout)
	var arr []string
	r, err := regexp.Compile(`^.{110}(/.*?)$`)
	if err != nil {
		log.Fatal(err)
	}
	for outScanner.Scan() {
		res := r.FindAllStringSubmatch(outScanner.Text(), -1)
		if res != nil {
			arr = append(arr, res[0][1])

		}
	}
	return arr
}
