/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"log"
	"os"

	"github.com/jim380/Cendermint/cmd"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/rest"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	chainList                                        = []string{"cosmos", "umee", "nym"}
	chain, restAddr, listenPort, operAddr, logOutput string
	logger                                           *zap.Logger
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chain = os.Getenv("CHAIN")
	operAddr = os.Getenv("OPERATOR_ADDR")
	restAddr = os.Getenv("REST_ADDR")
	listenPort = os.Getenv("LISTENING_PORT")
	logOutput = os.Getenv("LOG_OUTPUT")

	logger = logging.InitLogger(logOutput)
	zap.ReplaceGlobals(logger)

	cmd.CheckInputs(chain, operAddr, restAddr, listenPort, chainList)
	cmd.SetSDKConfig(chain)
	rest.Addr = restAddr
	rest.OperAddr = operAddr
	startExporter()
}

func startExporter() {
	exporter.Start(chain, listenPort, logger)
}
