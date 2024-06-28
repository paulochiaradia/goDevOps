package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	ServerName    string
	ServerUrl     string
	TempoExecucao float64
	StatusCode    int
	DownTime      string
}

func criarListaServidores(serverList *os.File) []Server {
	csvReader := csv.NewReader(serverList)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var servidores []Server

	for i, line := range data {
		if i > 0 {
			servidor := Server{
				ServerName: line[0],
				ServerUrl:  line[1],
			}
			servidores = append(servidores, servidor)
		}
	}

	return servidores
}

func checkServidores(servidores []Server) []Server {
	var downServers []Server
	for _, server := range servidores {
		agora := time.Now()
		get, err := http.Get(server.ServerUrl)
		if err != nil {
			fmt.Printf("Server %s is down [%s]\n", server.ServerName, err.Error())
			server.StatusCode = 0
			server.DownTime = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
			continue
		}
		server.TempoExecucao = time.Since(agora).Seconds()
		server.StatusCode = get.StatusCode
		if server.StatusCode != 200 {
			server.DownTime = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
		}
		fmt.Printf("ServerUrl: [%s]\nStatus: [%d] Tempo:[%f]\n", server.ServerUrl, server.StatusCode, server.TempoExecucao)
	}
	return downServers
}

func openFiles(serverListFile string, downTimeFile string) (*os.File, *os.File) {
	serverList, err := os.OpenFile(serverListFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	downTimeList, err := os.OpenFile(downTimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return serverList, downTimeList
}

func generateDownTime(downTimeList *os.File, downServers []Server) {
	csvWriter := csv.NewWriter(downTimeList)
	for _, servidor := range downServers {
		line := []string{servidor.ServerName, servidor.ServerUrl, servidor.DownTime, fmt.Sprintf("%f", servidor.TempoExecucao), fmt.Sprintf("%d", servidor.StatusCode)}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func main() {

	serverList, downTimeList := openFiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downTimeList.Close()
	servidores := criarListaServidores(serverList)
	downServers := checkServidores(servidores)
	generateDownTime(downTimeList, downServers)

}
