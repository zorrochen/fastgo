package main





type Server struct {
	pawObj	*parse_and_write_t
	exitChan chan bool
}

var global_server_ref *Server

func NewServer() *Server {
	return &Server{
		pawObj: new(parse_and_write_t),
	}
}

func (s *Server) Start() {
	global_server_ref = s

	s.pawObj.Parse()
}


func (s *Server) Exit() {

}

