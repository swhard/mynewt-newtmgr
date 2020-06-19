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

 package xact

 import (
	 "mynewt.apache.org/newtmgr/nmxact/nmp"
	 "mynewt.apache.org/newtmgr/nmxact/sesn"
 )
 
 type TestPepCmd struct {
	 CmdBase
	 Payload string
 }
 
 func NewTestPepCmd() *TestPepCmd {
	 return &TestPepCmd{
		 CmdBase: NewCmdBase(),
	 }
 }
 
 type TestPepResult struct {
	 Rsp *nmp.TestPepRsp
 }
 
 func newTestPepResult() *TestPepResult {
	 return &TestPepResult{}
 }
 
 func (r *TestPepResult) Status() int {
	 return r.Rsp.Rc
 }
 
 func (c *TestPepCmd) Run(s sesn.Sesn) (Result, error) {
	 r := nmp.NewTestPepReq()
	 r.Payload = c.Payload
 
	 rsp, err := txReq(s, r.Msg(), &c.CmdBase)
	 if err != nil {
		 return nil, err
	 }
	 srsp := rsp.(*nmp.TestPepRsp)
 
	 res := newTestPepResult()
	 res.Rsp = srsp
	 return res, nil
 }
 