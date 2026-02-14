# FCGI-Healthcheck

## Why does this exist?

`fgci-bin`, in a Kubernetes `livenessProbe`, was not behaving itself.

This led to `php-fpm` pods being terminated.

An alternative was desired. The primary requirement being that it works.

## Minimum configuration

```shell
fcgi-healthcheck --uri /ping
```

Just `/ping` for the URI is enough. To see the `pong` coming back, add the verbose flag.

If it worked, it exits with 0, if it failed, it exits with 1.

That should be sufficient for your healthchecks.

## Testing

A docker compose `php-fpm` setup is included here. `docker compose up -d` will get you started.

## Isn't this really basic?

Yes. The heavy lifting is being done by the library `github.com/tomasen/fcgi_client`.

## What else can we do?

You can execute PHP like so:

```shell
fcgi-healthcheck --script info.php -V
```

## Are you absolutely sure this needed to exist?

Why are you asking me this?
