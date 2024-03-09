package middleware

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// This function initializes New Relic and returns the application instance.
// Usually, you would call this function once and pass the application instance to your middleware.
func InitNewRelic(appName, licenseKey string) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func NewRelicMiddleware(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.URL.Path)
			defer txn.End()

			// Attach transaction to the request context
			r = newrelic.RequestWithTransactionContext(r, txn)

			next.ServeHTTP(w, r)
		})
	}
}

// NewRelicMiddleware creates a middleware function that starts a New Relic transaction at the beginning of each HTTP request.
/*func NewRelicMiddleware(app *newrelic.Application, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := app.StartTransaction(r.URL.Path)
		defer txn.End()

		// Before calling the next handler, set the transaction in the request's context
		r = newrelic.RequestWithTransactionContext(r, txn)

		// Set up response writer to capture status code for New Relic transaction
		w = txn.SetWebResponse(w)

		// Call the next handler, which can now retrieve the transaction from the context
		next.ServeHTTP(w, r)

		// You can add attributes to the transaction based on the request or response
		txn.AddAttribute("requestMethod", r.Method)
	}
}*/
