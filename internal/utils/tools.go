// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"github.com/ghodss/yaml"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
)

const ClientSetKey = "clientsetKey"

type Prefixes struct {
	PrefixesList []string `json:"Prefixes"`
}

func MergeNotifyChannels(channels ...<-chan struct{}) <-chan struct{} {
	out := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(channels))
	output := func(c <-chan struct{}) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FromContext(ctx context.Context) *kubernetes.Clientset {
	return ctx.Value(ClientSetKey).(*kubernetes.Clientset)
}

func PrefixesToYaml(prefixesList []string) ([]byte, error) {
	source := Prefixes{prefixesList}

	bytes, err := yaml.Marshal(source)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func YamlToPrefixes(bytes []byte) (Prefixes, error) {
	destination := Prefixes{}
	err := yaml.Unmarshal(bytes, &destination)
	if err != nil {
		return Prefixes{}, err
	}

	return destination, nil
}

func Notify(notifyChan chan<- struct{}) {
	notifyChan <- struct{}{}
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err == nil {
			return nil
		}
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
