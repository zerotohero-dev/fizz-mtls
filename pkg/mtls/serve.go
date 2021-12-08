/*
 *  \
 *  \\,
 *   \\\,^,.,,.                     Zero to Hero
 *   ,;7~((\))`;;,,               <zerotohero.dev>
 *   ,(@') ;)`))\;;',    stay up to date, be curious: learn
 *    )  . ),((  ))\;,
 *   /;`,,/7),)) )) )\,,
 *  (& )`   (,((,((;( ))\,
 */

package mtls

import (
	"context"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"github.com/zerotohero-dev/fizz-logging/pkg/log"
	"net"
)

type SpireArgs struct {
	AppPrefix      string
	AppNameDefault string
	AppName        string
	AppNameIdm     string
	AppNameMailer  string
	AppTrustDomain string
}

func AnyId(args SpireArgs) ([]spiffeid.ID, error) {
	anyId, err := spiffeid.New(
		args.AppTrustDomain, args.AppPrefix, args.AppNameDefault,
	)

	if err != nil {
		return nil, err
	}

	return []spiffeid.ID{anyId}, nil
}

func FromIdm(args SpireArgs) ([]spiffeid.ID, error) {
	appId, err := spiffeid.New(
		args.AppTrustDomain, args.AppPrefix, args.AppNameIdm,
	)
	if err != nil {
		return nil, err
	}

	return []spiffeid.ID{appId}, nil
}

func ListenAndServe(
	serverAddress, socketPath string,
	spireArgs SpireArgs,
	allowedIds []spiffeid.ID,
	svc func() interface{},
	handlerFn func(net.Conn, interface{}),
	errFn func(error),
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("SPIRE mTLS server will try listening… (%s)", serverAddress)

	listener, err := spiffetls.ListenWithMode(ctx, "tcp", serverAddress,
		spiffetls.MTLSServerWithSourceOptions(
			tlsconfig.AuthorizeOneOf(allowedIds...),
			workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)),
		))

	if err != nil {
		log.Err("SPIRE: Unable to create TLS listener: %v", err.Error())
		panic(err.Error())
	}

	log.Info(
		"SPIRE: 🐢 Service '%s' is waiting for mTLS connections at '%s",
		spireArgs.AppName, serverAddress,
	)

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Err("SPIRE: Possibly leaking a listener: '%s'", err.Error())
		}
	}()

	// svc := service.New(svcArgs, ctx)
	for {
		conn, err := listener.Accept()
		if err != nil {
			go errFn(err)
			continue
		}
		go handlerFn(conn, svc)
	}
}