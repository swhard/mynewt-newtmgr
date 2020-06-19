/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
	"mynewt.apache.org/newtmgr/nmxact/xact"
)

func testpepRunCmd(cmd *cobra.Command, args []string) {
	fmt.Println("PIPPO")
	if len(args) != 1 {
		nmUsage(cmd, nil)
	}

	s, err := GetSesn()
	if err != nil {
		nmUsage(nil, err)
	}

	c := xact.NewTestPepCmd()
	c.SetTxOptions(nmutil.TxOptions())
	c.Payload = args[0]

	res, err := c.Run(s)
	if err != nil {
		nmUsage(nil, util.ChildNewtError(err))
	}

	eres := res.(*xact.TestPepResult)
	fmt.Println(eres.Rsp.Payload)
}

func testPepCmd() *cobra.Command {
	testPepCmd := &cobra.Command{
		Use:   "testPep <text> -c <conn_profile>",
		Short: "Send data to a device and display the echoed back data",
		Run:   testpepRunCmd,
	}

	return testPepCmd
}
