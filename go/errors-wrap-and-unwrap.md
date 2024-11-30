# Errors wrap and unwrap
Should use `%w` instead of `%s` or `%v`, so that `errors.Unwrap` can unwrap the message and compare it using `errors.Is`.