package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/result"
	"github.com/spf13/cobra"
)

// runCmdFunc wraps cmdutil.RunFunc. While cmdutil.RunFunc provides a standard
// wrapper for dealing with and logging errors before exiting with an
// appropriate error code, runCmdFunc extends this with additional error
// handling specific to the Pulumi CLI. This includes e.g. specific and more
// helpful messages in the case of decryption or snapshot integrity errors.
func runCmdFunc(
	run func(cmd *cobra.Command, args []string) error,
) func(cmd *cobra.Command, args []string) {
	return cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
		err := run(cmd, args)
		return processCmdErrors(err)
	})
}

// Processes errors that may be returned from commands, providing a central
// location to insert more human-friendly messages when certain errors occur, or
// to perform other type-specific handling.
func processCmdErrors(err error) error {
	// If no error occurred, we have nothing to do.
	if err == nil {
		return nil
	}

	// If the error is a "bail" (that is, some expected error flow), then a
	// diagnostic or message will already have been reported. We can thus return
	// it in order to effect the exit code without printing the error message
	// again.
	if result.IsBail(err) {
		return err
	}

	// Other type-specific error handling.
	if de, ok := engine.AsDecryptError(err); ok {
		printDecryptError(*de)
		return nil
	}

	// In all other cases, return the unexpected error as-is for generic handling.
	return err
}

// A type-specific handler for engine.DecryptErrors that prints out help text
// containing common causes of and possible resolutions for decryption errors.
func printDecryptError(e engine.DecryptError) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	fprintf(writer, "failed to decrypt encrypted configuration value '%s': %s", e.Key, e.Err)
	fprintf(writer, ""+
		"This can occur when a secret is copied from one stack to another. Encryption of secrets is done per-stack and "+
		"it is not possible to share an encrypted configuration value across stacks.\n"+
		"\n"+
		"You can re-encrypt your configuration by running `pulumi config set %s [value] --secret` with your "+
		"new stack selected.\n"+
		"\n"+
		"refusing to proceed", e.Key)
	contract.IgnoreError(writer.Flush())
	cmdutil.Diag().Errorf(diag.RawMessage("" /*urn*/, buf.String()))
}

// Quick and dirty utility function for printing to writers that we know will never fail.
func fprintf(writer io.Writer, msg string, args ...interface{}) {
	_, err := fmt.Fprintf(writer, msg, args...)
	contract.IgnoreError(err)
}
