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

 package nmp

 import ()
 
 type TestPepReq struct {
	 NmpBase        `codec:"-"`
	 Payload string `codec:"d"`
 }
 
 type TestPepRsp struct {
	 NmpBase
	 Payload string `codec:"r"`
	 Rc      int    `codec:"rc"`
 }
 
 func NewTestPepReq() *TestPepReq {
	 r := &TestPepReq{}
	 fillNmpReq(r, NMP_OP_WRITE, NMP_GROUP_PEPP, NMP_ID_DEF_ECHO)
	 return r
 }
 
 func (r *TestPepReq) Msg() *NmpMsg { return MsgFromReq(r) }
 
 func NewTestPepRsp() *TestPepRsp {
	 return &TestPepRsp{}
 }
 
 func (r *TestPepRsp) Msg() *NmpMsg { return MsgFromReq(r) }
 