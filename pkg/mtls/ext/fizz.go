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

package ext

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/zerotohero-dev/fizz-entity/pkg/reqres"
	"net"
)

// FizzBuzz.Pro Extensions

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

func Idm(args SpireArgs) ([]spiffeid.ID, error) {
	appId, err := spiffeid.New(
		args.AppTrustDomain, args.AppPrefix, args.AppNameIdm,
	)
	if err != nil {
		return nil, err
	}

	return []spiffeid.ID{appId}, nil
}

func AllowList(args SpireArgs, allowAll bool) ([]spiffeid.ID, error) {
	if allowAll {
		res, err := AnyId(args)
		if err != nil {
			return []spiffeid.ID{}, errors.Wrap(err, "problem generating allow list")
		}
		return res, nil
	}

	res, err := Idm(args)
	if err != nil {
		return []spiffeid.ID{}, errors.Wrap(err, "problem generating allow list")
	}
	return res, nil
}

func Payload(request string) (string, error) {
	var req reqres.MtlsApiRequest

	err := json.Unmarshal([]byte(request), &req)
	if err != nil {
		return "", errors.Wrap(err, "payload: problem with unmarshal")
	}

	body := req.Body
	return body, nil
}

func Send(conn net.Conn, result interface{}) error {
	serialized, _ := json.Marshal(result)
	if _, err := conn.Write([]byte(string(serialized) + "\n")); err != nil {
		return errors.Wrap(err, "Unable to send a response")
	}

	return nil
}
