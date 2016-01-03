package main

type Server struct {
	pawObj   *parse_and_write_t
	tranFuncObj   *tran_module_t
	tranSqlObj   *tran_sql_t
	exitChan chan bool
}

var global_server_ref *Server

func NewServer() *Server {
	return &Server{
		pawObj: new(parse_and_write_t),
		tranFuncObj: new(tran_module_t),
		tranSqlObj: new(tran_sql_t),
	}
}

func (s *Server) Start() {
	global_server_ref = s

	s.pawObj.Parse()
}

func (s *Server) Exit() {

}
