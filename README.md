# Status Code/Error Replacer
[![Code Coverage](https://codecov.io/gh/PseudoResonance/traefikerrorreplace/branch/master/graph/badge.svg?token=QFGZS5QJSG)](https://codecov.io/gh/PseudoResonance/traefikerrorreplace)
[![Code Analysis](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/codeqlAnalysis.yml/badge.svg)](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/codeqlAnalysis.yml)
[![Codacy Security Scan](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/codacyAnalysis.yml/badge.svg)](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/codacyAnalysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/PseudoResonance/traefikerrorreplace)](https://goreportcard.com/report/github.com/PseudoResonance/traefikerrorreplace)
[![Build and Test Source](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/buildAndTest.yml/badge.svg)](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/buildAndTest.yml)
[![Static Analysis](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/staticAnalysis.yml/badge.svg)](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/staticAnalysis.yml)
[![Integration Test](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/prodTest.yml/badge.svg)](https://github.com/PseudoResonance/traefikerrorreplace/actions/workflows/prodTest.yml)

Some apps become unhappy when they receive unexpected error codes, such as when Traefik can't find an available server.

This plugin solves this issue by filtering and replacing problematic error codes with ones the app won't complain about.

## Configuration

### Configuration documentation

Supported configurations per body

| Setting       | Allowed values | Required | Description                   |
| :------------ | :------------- | :------- | :---------------------------- |
| matchStatus   | []int          | Yes      | List of status codes to match |
| replaceStatus | int            | Yes      | Status code to replace with   |
| debug         | bool           | No       | Enables extra debug logging   |

### Enable the plugin

```yaml
experimental:
  plugins:
    traefikerrorreplace:
      modulename: github.com/PseudoResonance/traefikerrorreplace
      version: v1.0.1
```

### Plugin configuration

```yaml
http:
  middlewares:
    traefikerrorreplace:
      plugin:
        traefikerrorreplace:
          matchStatus:
            - 500
            - 503
          replaceStatus: 404

  routers:
    my-router:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - traefikerrorreplace

  services:
    service-whoami:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
```

# Testing

[https://github.com/PseudoResonance/traefikerrorreplace/tree/master/test](https://github.com/PseudoResonance/traefikerrorreplace/tree/master/test)

We have written the following tests in this repo:

- golang linting
- yaegi tests (validate configuration matches what Traefik expects)
- General GO code coverage
- Virtual implementation tests (spin up traefik with yml/toml tests to make sure the plugin actually works)
- Live implementation tests (spin up traefik with the plugin definition as it would be for you, and run the same tests again)

These tests allow us to make sure the plugin is always functional with Traefik and Traefik version updates.
