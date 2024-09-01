// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package record

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

type LogRecorder struct {
	Logger        logr.Logger
	IncludeObject bool
}

func (l *LogRecorder) logger() logr.Logger {
	if l.Logger.GetSink() == nil {
		return logr.Discard()
	} else {
		return l.Logger
	}
}

func (l *LogRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	if l.IncludeObject {
		l.logger().Info(message, "reason", reason, "eventType", eventtype, "object", object)
	} else {
		l.logger().Info(message, "reason", reason, "eventType", eventtype)
	}
}

func (l *LogRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	msg := fmt.Sprintf(messageFmt, args...)
	if l.IncludeObject {
		l.logger().Info(msg, "reason", reason, "eventType", eventtype, "object", object)
	} else {
		l.logger().Info(msg, "reason", reason, "eventType", eventtype)
	}
}

func (l *LogRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	msg := fmt.Sprintf(messageFmt, args...)
	if l.IncludeObject {
		l.logger().Info(msg, "reason", reason, "eventType", eventtype, "annotations", annotations, "object", object)
	} else {
		l.logger().Info(msg, "reason", reason, "eventType", eventtype, "annotations", annotations)
	}
}
