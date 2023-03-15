// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetryapireceiver // import "github.com/open-telemetry/opentelemetry-lambda/collector/receiver/telemetryapireceiver"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-lambda/collector/internal/sharedcomponent"
	//"go.opentelemetry.io/collector/internal/sharedcomponent"
)

const (
	typeStr   = "telemetryapi"
	stability = component.StabilityLevelDevelopment
)

var errConfigNotTelemetryAPI = errors.New("config was not a Telemetry API receiver config")

// NewFactory creates a new receiver factory
func NewFactory(extensionID string) receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		func() component.Config {
			return &Config{
				extensionID: extensionID,
			}
		},
		receiver.WithTraces(createTracesReceiver, stability),
		receiver.WithLogs(createLogsReceiver, stability))
}

// createTraces creates a trace receiver based on provided config.
func createTracesReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Traces,
) (receiver.Traces, error) {
	oCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errConfigNotTelemetryAPI
	}
	r, err := receivers.GetOrAdd(oCfg, func() (*telemetryAPIReceiver, error) {
		return newTelemetryAPIReceiver(oCfg, set)
	})
	if err != nil {
		return nil, err
	}

	//r.tracesConsumer = consumer

	return r, nil
}

// createLogsReceiver creates a logs receiver based on provided config.
func createLogsReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	oCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errConfigNotTelemetryAPI
	}
	r, err := receivers.GetOrAdd(oCfg, func() (*telemetryAPIReceiver, error) {
		return newTelemetryAPIReceiver(oCfg, set)
	})
	if err != nil {
		return nil, err
	}
	r.Unwrap().registerLogsConsumer(consumer)

	return r, nil
}

// https://github.com/open-telemetry/opentelemetry-collector/blob/v0.73.0/receiver/otlpreceiver/factory.go#L134
// This is the map of already created receivers for particular configurations.
// We maintain this map because the Factory is asked trace and metric receivers separately
// when it gets CreateTracesReceiver() and CreateMetricsReceiver() but they must not
// create separate objects, they must use one Receiver object per configuration.
// When the receiver is shutdown it should be removed from this map so the same configuration
// can be recreated successfully.
var receivers = sharedcomponent.NewSharedComponents[*Config, *telemetryAPIReceiver]()
