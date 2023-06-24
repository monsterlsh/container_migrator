package client

// TODO 将一些功能集合成struct

type Socket struct {
	client_path string
	server_path string
	containerID string
	destIP      string
	othersPath  string
	destPath    string
	basePath    string
}

func (s *Socket) SetBasePath() {

}
