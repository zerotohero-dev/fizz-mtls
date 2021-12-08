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
	"github.com/pkg/errors"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/zerotohero-dev/fizz-mtls/pkg/mtls"
)

// FizzBuzz.Pro Extensions

func AnyId(args mtls.SpireArgs) ([]spiffeid.ID, error) {
	anyId, err := spiffeid.New(
		args.AppTrustDomain, args.AppPrefix, args.AppNameDefault,
	)

	if err != nil {
		return nil, err
	}

	return []spiffeid.ID{anyId}, nil
}

func Idm(args mtls.SpireArgs) ([]spiffeid.ID, error) {
	appId, err := spiffeid.New(
		args.AppTrustDomain, args.AppPrefix, args.AppNameIdm,
	)
	if err != nil {
		return nil, err
	}

	return []spiffeid.ID{appId}, nil
}

func AllowList(args mtls.SpireArgs, allowAll bool) ([]spiffeid.ID, error) {
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