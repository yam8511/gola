package def

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

// HttpTraceInfo 追蹤資訊
type HttpTraceInfo struct {
	Begin, End,
	DnsStart,
	DnsDone,
	TcpConnStart,
	TcpConnDone,
	HandshakeStart,
	HandshakeDone,
	ServerProcessStart,
	ServerProcessDone,
	TransferStart,
	TransferDone,
	ResultParseStart,
	ResultParseDone,
	BodyCloseStart,
	BodyCloseDone time.Time
	BeginHost,
	SelectedIP string
	AvailableIP []string
	TraceErr    error
}

// CalculatedHttpTraceInfo 計算追蹤資訊
type CalculatedHttpTraceInfo struct {
	Begin, End time.Time
	DnsLookup,
	TcpConnection,
	TlsHandshake,
	ServerProcessing,
	ContentTransfer,
	ResultParse,
	BodyClose,
	NameLookup,
	Connect,
	PreTransfer,
	StartTransfer,
	Total string
	BeginHost,
	SelectedIP string
	AvailableIP []string
	TraceErr    error
}

// Calculate 計算追蹤時間
//   DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer    Result Parse       Body Close
// [     24ms  |          42ms  |        371ms  |            354ms  |             0ms  |           0ms   |            0ms  ]
//             |                |               |                   |                  |                 |                 |
//    namelookup:24ms           |               |                   |                  |                 |                 |
//                        connect:66ms          |                   |                  |                 |                 |
//                                    pretransfer:437ms             |                  |                 |                 |
//                                                      starttransfer:792ms            |                 |                 |
//                                                                               receive:792ms           |                 |
//                                                                                                   parse:792ms           |
//                                                                                                                    total:792ms
func (trace *HttpTraceInfo) Calculate() *CalculatedHttpTraceInfo {
	if trace == nil {
		return nil
	}
	calTrace := CalculatedHttpTraceInfo{
		Begin:       trace.Begin,
		End:         trace.End,
		BeginHost:   trace.BeginHost,
		AvailableIP: trace.AvailableIP,
		SelectedIP:  trace.SelectedIP,
		TraceErr:    trace.TraceErr,
	}

	calTrace.DnsLookup = trace.DnsDone.Sub(trace.DnsStart).String()
	calTrace.NameLookup = trace.TcpConnStart.Sub(trace.Begin).String()
	calTrace.TcpConnection = trace.TcpConnDone.Sub(trace.TcpConnStart).String()
	calTrace.Connect = trace.HandshakeStart.Sub(trace.Begin).String()
	calTrace.TlsHandshake = trace.HandshakeDone.Sub(trace.HandshakeStart).String()
	calTrace.PreTransfer = trace.ServerProcessStart.Sub(trace.Begin).String()
	calTrace.ServerProcessing = trace.ServerProcessDone.Sub(trace.ServerProcessStart).String()
	calTrace.ServerProcessing = "<header>(" + calTrace.ServerProcessing + ") + <body>(" + trace.TransferStart.Sub(trace.ServerProcessDone).String() + ")"
	calTrace.StartTransfer = trace.TransferStart.Sub(trace.Begin).String()
	calTrace.ContentTransfer = trace.TransferDone.Sub(trace.TransferStart).String()
	calTrace.ResultParse = trace.ResultParseDone.Sub(trace.ResultParseStart).String()
	calTrace.BodyClose = trace.BodyCloseDone.Sub(trace.BodyCloseStart).String()
	calTrace.Total = trace.End.Sub(trace.Begin).String()

	return &calTrace
}

func (info *CalculatedHttpTraceInfo) String() string {
	var msg string
	msg += fmt.Sprintln("\n可選IP： ", info.AvailableIP)
	msg += fmt.Sprintln("選定IP： ", info.SelectedIP)
	msg += fmt.Sprintln("連線host： ", info.BeginHost)
	msg += fmt.Sprintln("追蹤錯誤： ", info.TraceErr)
	msg += fmt.Sprintln("================================")
	msg += fmt.Sprintf("* Begin 【 %s 】\n", info.Begin)
	if info.DnsLookup != "" {
		msg += fmt.Sprintf("    + DNS Lookup        - %s\n", info.DnsLookup)
	}
	if info.NameLookup != "" {
		msg += fmt.Sprintf("* namelookup:           - %s\n", info.NameLookup)
	}
	if info.TcpConnection != "" {
		msg += fmt.Sprintf("    + TCP Connection    - %s\n", info.TcpConnection)
	}
	if info.Connect != "" {
		msg += fmt.Sprintf("* connect:              - %s\n", info.Connect)
	}
	if info.TlsHandshake != "" {
		msg += fmt.Sprintf("    + TLS Handshake     - %s\n", info.TlsHandshake)
	}
	if info.PreTransfer != "" {
		msg += fmt.Sprintf("* pretransfer:          - %s\n", info.PreTransfer)
	}
	if info.ServerProcessing != "" {
		msg += fmt.Sprintf("    + Server Processing - %s\n", info.ServerProcessing)
	}
	if info.StartTransfer != "" {
		msg += fmt.Sprintf("* starttransfer:        - %s\n", info.StartTransfer)
	}
	if info.ContentTransfer != "" {
		msg += fmt.Sprintf("    + Content Transfer  - %s\n", info.ContentTransfer)
	}
	if info.ResultParse != "" {
		msg += fmt.Sprintf("    + Result Parse      - %s\n", info.ResultParse)
	}
	if info.BodyClose != "" {
		msg += fmt.Sprintf("    + Body Close        - %s\n", info.BodyClose)
	}
	msg += fmt.Sprintf("* End 【 %s 】\n", info.End)
	msg += fmt.Sprintf("    + total:            - %s\n", info.Total)

	return msg
}

// WithHttpTrace 替請求加上追蹤功能
func WithHttpTrace(req *http.Request) (*http.Request, *HttpTraceInfo) {

	traceInfo := &HttpTraceInfo{}

	// 追蹤線路時間
	trace := &httptrace.ClientTrace{
		// 開始
		GetConn: func(hostPort string) {
			traceInfo.BeginHost = hostPort
		},

		// DNS Lookup (namelookup)
		DNSStart: func(httptrace.DNSStartInfo) {
			if traceInfo.DnsStart.IsZero() {
				traceInfo.DnsStart = time.Now()
			}
		},
		DNSDone: func(doneInfo httptrace.DNSDoneInfo) {
			traceInfo.DnsDone = time.Now()
			if doneInfo.Err != nil {
				traceInfo.TraceErr = doneInfo.Err
			} else {
				for _, ip := range doneInfo.Addrs {
					traceInfo.AvailableIP = append(traceInfo.AvailableIP, ip.IP.String())
				}
			}
		},

		// TCP Connection (connect)
		ConnectStart: func(network, addr string) {
			if traceInfo.TcpConnStart.IsZero() {
				traceInfo.TcpConnStart = time.Now()
			}
		},
		ConnectDone: func(network, addr string, err error) {
			traceInfo.TcpConnDone = time.Now()

			if err != nil {
				traceInfo.TraceErr = err
			} else {
				traceInfo.SelectedIP = addr
			}
		},

		// TLS Handshake (pretransfer)
		TLSHandshakeStart: func() {
			if traceInfo.HandshakeStart.IsZero() {
				traceInfo.HandshakeStart = time.Now()
			}
		},
		TLSHandshakeDone: func(info tls.ConnectionState, err error) {
			traceInfo.HandshakeDone = time.Now()
			if err != nil {
				traceInfo.TraceErr = err
			}
		},

		// Server Processing (starttransfer)
		GotConn: func(info httptrace.GotConnInfo) {
			if traceInfo.ServerProcessStart.IsZero() {
				traceInfo.ServerProcessStart = time.Now()
			}
		},

		WroteRequest: func(info httptrace.WroteRequestInfo) {
			traceInfo.ServerProcessDone = time.Now()

			if info.Err != nil {
				traceInfo.TraceErr = info.Err
			}
		},

		// Content Transfer
		GotFirstResponseByte: func() {
			if traceInfo.TransferStart.IsZero() {
				traceInfo.TransferStart = time.Now()
			}
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace)), traceInfo
}
