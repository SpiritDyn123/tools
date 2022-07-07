package main

import "flag"

var (
	host string
	port string
	show_cmd  bool
	num_to_str bool
)


func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "redis host")
	flag.StringVar(&port, "p", "6300", "redis port")
	flag.BoolVar(&show_cmd, "sc", false, "print cmd")
	flag.BoolVar(&num_to_str, "nts", true, "inter number to string")
	flag.Parse()
}

