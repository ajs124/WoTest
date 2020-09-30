package runner

import (
	"context"
	"errors"
	"github.com/ajs124/WoTest/clients"
	. "github.com/ajs124/WoTest/config"
	"github.com/rs/zerolog/log"
	"math"
	"net/url"
	"sort"
	"sync"
	"time"
)

func asyncRecv(client clients.Client, reqUrl *url.URL, reqNum int, mdC chan clients.MeasurementData, wg *sync.WaitGroup, last bool) {
	_, md, err := client.Recv(*reqUrl)
	if err != nil {
		log.Err(err).Int("request number", reqNum).Msg("http client error")
	} else {
		if reqNum%32 == 0 {
			log.Debug().Int("request number", reqNum).Msg("succeeded")
		}
	}
	mdC <- md
	wg.Done()
	if last { // don't tell anyone you saw this
		wg.Wait()
		close(mdC)
	}
}

func getTotal(rm []clients.MeasurementData) time.Duration {
	sum := 0 * time.Second
	for _, m := range rm {
		sum += m.StopTime.Sub(m.StartTime)
	}
	return sum
}

func getPercentile(rm []clients.MeasurementData, percentile int) time.Duration {
	times := make([]time.Duration, 0)
	for _, m := range rm {
		dur := m.StopTime.Sub(m.StartTime)
		times = append(times, dur)
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	var idx float64
	idx = float64(percentile*len(times)) / 100
	return times[int(math.Round(idx))]
}

func getMean(rm []clients.MeasurementData) time.Duration {
	total := getTotal(rm)
	return time.Nanosecond * time.Duration(total.Nanoseconds()/int64(len(rm)))
}

func getStdDev(rm []clients.MeasurementData) time.Duration {
	sum := 0.0
	avg := float64(getMean(rm))
	for _, m := range rm {
		dur := float64(m.StopTime.Sub(m.StartTime))
		sum += math.Pow(dur-avg, 2)
	}
	fStdDev := math.Sqrt(sum / float64(len(rm)))
	return time.Duration(math.Round(fStdDev))
}

func runMeasureTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	// set up data structures etc
	var result TestResult
	result.measurements = make(map[int]*TestMeasurements)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(test.Timeout))
	stdout := make([]byte, 0)
	stderr := make([]byte, 0)
	// run command
	cmd, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, &stdout, &stderr, append([]string{test.Path}, test.Args...)...)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}

	// wait for things to start up
	matches := 0
	for _, mM := range test.MeasureTestProperties.MustMatch {
		found := result.waitForOutput(ctx, &stdout, &stderr, mM)
		log.Debug().Str("mM", mM).Bool("found", found).Msg("found")
		if found {
			matches += 1
		}
	}
	log.Debug().Int("matches", matches).Int("out of", len(test.MeasureTestProperties.MustMatch)).Msg("")
	reqUrl, _ := url.Parse(test.MeasureTestProperties.RequestUrl)

	var client clients.Client

	if test.MeasureTestProperties.Protocol == ProtoHttp || test.MeasureTestProperties.Protocol == ProtoHttps {
		client = &clients.HttpClient{AllowSelfSignedCerts: true}
	} else if test.MeasureTestProperties.Protocol == ProtoCoap {
		client = &clients.CoapClient{}
	} else if test.MeasureTestProperties.Protocol == ProtoMqtt {
		return result, errors.New("not implemented")
	} else {
		// wtf?
		panic(errors.New("test has invalid protocol"))
	}

	for setNum, rs := range test.MeasureTestProperties.RequestSets {
		_, err = client.Connect(*reqUrl)
		if err != nil {
			log.Info().Err(err).Msg("error connecting")
			break
		}
		p := 0
		mdC := make(chan clients.MeasurementData)
		var wg sync.WaitGroup
		var measures TestMeasurements
		// start initial pool of workers
		for p < rs.Parallel {
			wg.Add(1)
			go asyncRecv(client, reqUrl, p, mdC, &wg, false)
			p += 1
		}
		reqNum := p
		for reqNum < rs.Num {
			// wait for result
			measures.raw = append(measures.raw, <-mdC)

			// start a new worker for every result
			wg.Add(1)
			go asyncRecv(client, reqUrl, reqNum, mdC, &wg, reqNum == rs.Num-1)
			reqNum += 1
		}
		// collect rest
		for {
			md, ok := <-mdC
			if ok {
				measures.raw = append(measures.raw, md)
			} else {
				break
			}
		}
		result.succeeded = len(test.MeasureTestProperties.MustMatch) == matches
		measures.total = getTotal(measures.raw)
		measures.mean = getMean(measures.raw)
		measures.median = getPercentile(measures.raw, 50)
		measures.nfth = getPercentile(measures.raw, 95)
		measures.nnth = getPercentile(measures.raw, 99)
		measures.stddev = getStdDev(measures.raw)
		result.measurements[setNum] = &measures
		err = client.Disconnect()
	}

	result.stdout = string(stdout)
	result.stderr = string(stderr)
	cancel()
	err = cmd.Process.Kill()

	return result, err
}
