// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/Ungigdu/BAS_contract_go/BAS_Ethereum"
	"github.com/kprc/basserver/dns/server"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"

	"github.com/kprc/basserver/app/cmdcommon"
	"github.com/kprc/basserver/app/cmdservice"
	"github.com/kprc/basserver/config"
	"path"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "basd start in backend",
	Long:  `basd start in backend`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()
		config.GetBasDCfg().Save()

		daemondir := config.GetBASDHomeDir()
		cntxt := daemon.Context{
			PidFileName: path.Join(daemondir, "basd.pid"),
			PidFilePerm: 0644,
			LogFileName: path.Join(daemondir, "basd.log"),
			LogFilePerm: 0640,
			WorkDir:     daemondir,
			Umask:       027,
			Args:        []string{},
		}
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			log.Println("basd starting, please check log at:", path.Join(daemondir, "basd.log"))
			return
		}
		defer cntxt.Release()

		BAS_Ethereum.RecoverContract()
		go server.DNSServerDaemon()
		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	daemonCmd.Flags().IntVarP(&cmdroottcpport, "tcp-listen-port", "t", 65566, "local tcp listen port")
	daemonCmd.Flags().IntVarP(&cmdrootudpport, "udp-listen-port", "u", 65566, "local udp listen port")
	daemonCmd.Flags().StringVarP(&cmdropstennap, "ropsten-network-access-point", "r", "", "ropsten network access point")
	daemonCmd.Flags().StringVarP(&cmdbastokenaddr, "bas-token-address", "a", "", "bas token address")
	daemonCmd.Flags().StringVarP(&cmdbasmgraddr, "bas-mgr-address", "m", "", "bas manager address")
	daemonCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")

}
