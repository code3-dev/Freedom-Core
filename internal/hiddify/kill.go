package hiddify

import "context"

func KillHiddifyCrossPlatform(ctx context.Context) error {
	return KillHiddify(ctx)
}
