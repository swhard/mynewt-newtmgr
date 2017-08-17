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

package sesn

import (
	"github.com/runtimeco/go-coap"

	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/oic"
)

func TxNmp(s Sesn, m *nmp.NmpMsg, o TxOptions) (nmp.NmpRsp, error) {
	retries := o.Tries - 1
	for i := 0; ; i++ {
		r, err := s.TxNmpOnce(m, o)
		if err == nil {
			return r, nil
		}

		if !nmxutil.IsRspTimeout(err) || i >= retries {
			return nil, err
		}
	}
}

func getResourceOnce(s Sesn, resType ResourceType,
	uri string, opt TxOptions) (coap.COAPCode, []byte, error) {

	req, err := oic.CreateGet(true, uri, nmxutil.NextToken())
	if err != nil {
		return 0, nil, err
	}

	return s.TxCoapOnce(req, resType, opt)
}

func putResourceOnce(s Sesn, resType ResourceType,
	uri string, value []byte,
	opt TxOptions) (coap.COAPCode, []byte, error) {

	req, err := oic.CreatePut(true, uri, nmxutil.NextToken(), value)
	if err != nil {
		return 0, nil, err
	}

	return s.TxCoapOnce(req, resType, opt)
}

func postResourceOnce(s Sesn, resType ResourceType,
	uri string, value []byte,
	opt TxOptions) (coap.COAPCode, []byte, error) {

	req, err := oic.CreatePost(true, uri, nmxutil.NextToken(), value)
	if err != nil {
		return 0, nil, err
	}

	return s.TxCoapOnce(req, resType, opt)
}

func deleteResourceOnce(s Sesn, resType ResourceType,
	uri string, opt TxOptions) (coap.COAPCode, []byte, error) {

	req, err := oic.CreateDelete(true, uri, nmxutil.NextToken())
	if err != nil {
		return 0, nil, err
	}

	return s.TxCoapOnce(req, resType, opt)
}

func txCoap(txCb func() (coap.COAPCode, []byte, error),
	tries int) (coap.COAPCode, []byte, error) {

	retries := tries - 1
	for i := 0; ; i++ {
		code, r, err := txCb()
		if err == nil {
			return code, r, nil
		}

		if !nmxutil.IsRspTimeout(err) || i >= retries {
			return code, nil, err
		}
	}
}

func GetResource(s Sesn, resType ResourceType, uri string, o TxOptions) (
	coap.COAPCode, []byte, error) {

	return txCoap(func() (coap.COAPCode, []byte, error) {
		return getResourceOnce(s, resType, uri, o)
	}, o.Tries)
}

func PutResource(s Sesn, resType ResourceType, uri string,
	value []byte, o TxOptions) (coap.COAPCode, []byte, error) {

	return txCoap(func() (coap.COAPCode, []byte, error) {
		return putResourceOnce(s, resType, uri, value, o)
	}, o.Tries)
}

func PostResource(s Sesn, resType ResourceType, uri string,
	value []byte, o TxOptions) (coap.COAPCode, []byte, error) {

	return txCoap(func() (coap.COAPCode, []byte, error) {
		return postResourceOnce(s, resType, uri, value, o)
	}, o.Tries)
}

func DeleteResource(s Sesn, resType ResourceType, uri string,
	value []byte, o TxOptions) (coap.COAPCode, []byte, error) {

	return txCoap(func() (coap.COAPCode, []byte, error) {
		return deleteResourceOnce(s, resType, uri, o)
	}, o.Tries)
}

func PutCborResource(s Sesn, resType ResourceType, uri string,
	value map[string]interface{},
	o TxOptions) (coap.COAPCode, map[string]interface{}, error) {

	b, err := nmxutil.EncodeCborMap(value)
	if err != nil {
		return 0, nil, err
	}

	code, r, err := PutResource(s, resType, uri, b, o)
	m, err := nmxutil.DecodeCborMap(r)
	if err != nil {
		return 0, nil, err
	}

	return code, m, nil
}
