package odahu.core

import data.odahu.mapper
import data.odahu.roles

# Endpoints of sidecar containers
allow {
  mapper.action == "GET"
  mapper.resource == "/healthz"
}

allow {
  mapper.action == "GET"
  mapper.resource == "/metrics"
}

allow {
  mapper.action == "GET"
  mapper.resource == "/healthcheck"
}

# Triton Health endpoints
allow {
  mapper.action == "GET"
  mapper.resource == "/v2/health/live"
}

allow {
  mapper.action == "GET"
  mapper.resource == "/v2/health/ready"
}

# Full access for Admin
allow {
  mapper.user_roles[_] == roles.admin
}

# Triton Model endpoints
tritonPathRegex := `^/v2/models/[\w-]+(/versions/[\d]+)?(/infer|/ready)?/?$`

# Access to inference for data scientinsts and users with per-model role
allow {
  mapper.action == ["GET", "POST"][_]
  re_match(tritonPathRegex, mapper.resource)

  mapper.raw_roles[_] == ["{{.Role}}", roles.data_scientist][_]
}
